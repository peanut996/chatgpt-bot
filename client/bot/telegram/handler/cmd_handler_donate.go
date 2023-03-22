package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DonateCommandHandler struct{}

func (d *DonateCommandHandler) Cmd() BotCmd {
	return cmd.DONATE
}

func (d *DonateCommandHandler) Run(bot telegram.TelegramBot, message tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, tip.DonateTip)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.SafeSend(msg)

	return nil
}
func NewDonateCommandHandler() *DonateCommandHandler {
	return &DonateCommandHandler{}
}
