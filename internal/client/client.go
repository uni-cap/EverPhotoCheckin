package client

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/winterssy/ghttp"

	"github.com/winterssy/EverPhotoCheckin/internal/model"
)

const (
	_apiLogin   = "https://api.everphoto.cn/auth"
	_apiCheckin = "https://openapi.everphoto.cn/sf/3/v4/PostCheckIn"
)

var (
	_headers = ghttp.Headers{
		"user-agent": "EverPhoto 5.1.0 rv:5.1.0.0 (iPhone; iOS 15.6; zh_CN) Cronet",
	}
)

type Bot struct {
	client *ghttp.Client
}

const _salt = "tc.everphoto."

func salt(value string) string {
	hash := md5.Sum([]byte(_salt + value))
	return hex.EncodeToString(hash[:])
}

func New(mobile, password string) (*Bot, error) {
	client := ghttp.New()
	client.RegisterBeforeRequestCallbacks(
		ghttp.WithHeaders(_headers),
		ghttp.WithRetrier(),
	)

	resp, err := client.Post(_apiLogin, ghttp.WithForm(ghttp.Form{
		"mobile":   mobile,
		"password": salt(password),
	}))
	if err != nil {
		return nil, err
	}

	ar := new(model.AuthResponse)
	if err = resp.JSON(ar); err != nil {
		return nil, err
	}

	if ar.Code != 0 {
		return nil, errors.New(ar.Message)
	}

	client.RegisterBeforeRequestCallbacks(ghttp.WithBearerToken(ar.Data.Token))
	return &Bot{client: client}, nil
}

func NewWithToken(token string) *Bot {
	client := ghttp.New()
	client.RegisterBeforeRequestCallbacks(
		ghttp.WithHeaders(_headers),
		ghttp.WithBearerToken(token),
		ghttp.WithRetrier(),
	)
	return &Bot{client: client}
}

func (bot *Bot) Checkin(ctx context.Context) (*model.CheckinResult, error) {
	resp, err := bot.client.Post(_apiCheckin,
		ghttp.WithContentType("application/json"),
		ghttp.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	cr := new(model.CheckinResponse)
	if err = resp.JSON(cr); err != nil {
		return nil, err
	}

	if cr.Code != 0 {
		return nil, errors.New(cr.Message)
	}

	return cr.Data, nil
}
