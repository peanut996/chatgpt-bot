package limiter

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Limiter interface {
	Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string)
	CallBack(bot telegram.TelegramBot, message tgbotapi.Message, success bool)
}

func IsGPTMessage(message tgbotapi.Message) bool {
	return message.IsCommand() && (message.Command() == cmd.GPT4 || message.Command() == cmd.GPT)
}
