package telegram

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type TelegramBot interface {
	Init(*cfg.Config) error
	Run()

	SelfID() int64
	Config() *cfg.Config
	TGBot() *tgbotapi.BotAPI
	IsBotAdmin(from int64) bool

	GetBotInviteLink(code string) string
	GetUserInfo(userID int64) (*model.User, error)

	SafeSend(msg tgbotapi.MessageConfig)
	SafeReplyMsgWithoutPreview(chatID int64, messageID int, text string)
	SafeSendMsg(chatID int64, text string)
	SendAutoDeleteMessage(msg tgbotapi.MessageConfig, duration time.Duration)
}
