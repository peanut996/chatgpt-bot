package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type DonateCommandHandler struct{}

func (d *DonateCommandHandler) Cmd() BotCmd {
	return cmd.DONATE
}

func (d *DonateCommandHandler) Run(bot telegram.TelegramBot, message tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, tip.DonateTip)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.SendAutoDeleteMessage(msg, time.Second*30)
	return nil
}
func NewDonateCommandHandler() *DonateCommandHandler {
	return &DonateCommandHandler{}
}
