package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DowngradeCommandHandler struct {
}

func (c *DowngradeCommandHandler) Cmd() BotCmd {
	return cmd.DOWN
}

func (c *DowngradeCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	if !b.IsBotAdmin(message.From.ID) {
		msg.Text = tip.NotAdminTip
	} else {
		if message.CommandArguments() == "gpt" {
			b.Config().Downgrade = !b.Config().Downgrade
			msg.Text = fmt.Sprintf("downgrade mode is %v now", b.Config().Downgrade)
		}
		if message.CommandArguments() == "rate" {
			b.Config().EnableRateLimiter = !b.Config().EnableRateLimiter
			msg.Text = fmt.Sprintf("enable rate limiter is %v now", b.Config().EnableRateLimiter)
		}
		msg.Text = "Do nothing. Available is `gpt` or `rate`"
	}
	b.SafeSend(msg)
	return nil
}

func NewLimiterCommandHandler() *DowngradeCommandHandler {
	return &DowngradeCommandHandler{}
}
