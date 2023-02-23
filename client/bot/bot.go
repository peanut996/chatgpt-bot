package bot

import "chatgpt-bot/cfg"

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
		return NewTelegramBot()
	default:
		return NewTelegramBot()
	}
}
