package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/winterssy/EverPhotoCheckin/internal/client"
)

const (
	EnvEverPhotoMobile   = "EverPhotoMobile"
	EnvEverPhotoPassword = "EverPhotoPassword"
	EnvEverPhotoToken    = "EverPhotoToken"
)

var (
	_mobile   string
	_password string
	_token    string
)

func init() {
	log.SetPrefix("【时光相册】")
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	flag.StringVar(&_mobile, "mobile", "", "your mobile phone number")
	flag.StringVar(&_password, "password", "", "your password")
	flag.StringVar(&_token, "token", "", "your token")
}

func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

func createBot() (bot *client.Bot, err error) {
	_token = valueOrDefault(_token, os.Getenv(EnvEverPhotoToken))
	if _token != "" {
		bot = client.NewWithToken(_token)
		return
	}

	_mobile = valueOrDefault(_mobile, os.Getenv(EnvEverPhotoMobile))
	_password = valueOrDefault(_password, os.Getenv(EnvEverPhotoPassword))
	bot, err = client.New(_mobile, _password)
	return
}

func main() {
	flag.Parse()

	bot, err := createBot()
	if err != nil {
		log.Fatalf("登录失败：%s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cr, err := bot.Checkin(ctx)
	if err != nil {
		log.Fatalf("签到失败：%s" + err.Error())
	}

	log.Printf("你已连续签到%d天，累计获得空间%s，明天可白嫖%s，请继续保持(￣▽￣)", cr.Continuity, cr.TotalReward, cr.TomorrowReward)
}
