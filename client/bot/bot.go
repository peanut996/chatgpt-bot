package bot

import (
	"chatgpt-bot/bot/telegram/service"
	"chatgpt-bot/bot/wechat"
	"chatgpt-bot/cfg"
)

var (
	Telegram = "telegram"
	Wechat   = "wechat"
)

type Bot interface {
	Init(*cfg.Config) error
	Run()
}

func GetBot(botType string) Bot {
	switch botType {
	case Wechat:
		return wechat.NewWechatBot()
	case Telegram:
		return service.NewTelegramBot()
	default:
		return service.NewTelegramBot()
	}
}
