package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LimiterCommandHandler struct {
}

func (c *LimiterCommandHandler) Cmd() BotCmd {
	return cmd.LIMITER
}

func (c *LimiterCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	if !b.IsBotAdmin(message.From.ID) {
		msg.Text = tip.NotAdminTip
	} else {
		limiter := utils.ParseBoolString(message.CommandArguments())
		b.Config().RateLimiterConfig.Enable = limiter
		msg.Text = fmt.Sprintf("limiter status is %v now", limiter)
	}
	b.SafeSend(msg)
	return nil
}

func NewLimiterCommandHandler() *LimiterCommandHandler {
	return &LimiterCommandHandler{}
}
