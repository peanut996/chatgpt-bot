package bot

import "chatgpt-bot/cfg"

var (
	BOT_TELEGRAM = "telegram"
	BOT_WECHAT   = "wechat"
)

type Bot interface {
	Init(*cfg.Config) error
	Run()
}

func GetBot(botType string) Bot {
	switch botType {
	case BOT_WECHAT:
		return NewWechatBot()
	case BOT_TELEGRAM:
		return NewTelegramBot()
	default:
		return NewTelegramBot()
	}
}
