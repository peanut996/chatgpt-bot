package limiter

import (
	"chatgpt-bot/bot/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Limiter interface {
	Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string)
	CallBack(bot telegram.TelegramBot, message tgbotapi.Message, success bool)
}
