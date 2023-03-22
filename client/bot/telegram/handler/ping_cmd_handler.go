package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PingCommandHandler struct {
}

func (c *PingCommandHandler) Cmd() BotCmd {
	return cmd.PING
}

func (c *PingCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	b.SafeSendMsg(message.Chat.ID, tip.BotPingTip)
	return nil
}
func NewPingCommandHandler() *PingCommandHandler {
	return &PingCommandHandler{}
}
