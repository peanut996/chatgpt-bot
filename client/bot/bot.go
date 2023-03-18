package bot

import (
	"chatgpt-bot/bot/telegram"
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
		return NewWechatBot()
	case Telegram:
		return telegram.NewTelegramBot()
	default:
		return telegram.NewTelegramBot()
	}
}
