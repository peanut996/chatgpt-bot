package telegram

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot interface {
	Init(*cfg.Config) error
	Run()
	IsBotAdmin(from int64) bool
	SafeSend(msg tgbotapi.MessageConfig)
	GetUserInfo(userID int64) (*model.User, error)
	SafeReplyMsg(chatID int64, messageID int, text string)
	GetBotInviteLink(code string) string
	SafeSendMsg(chatID int64, text string)
	SelfID() int64
	Config() *cfg.Config
	GetAPIBot() *tgbotapi.BotAPI
}
