package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PushCommandHandler struct {
	userRepository *repository.UserRepository
}

func NewPushCommandHandler(userRepository *repository.UserRepository) *PushCommandHandler {
	return &PushCommandHandler{
		userRepository: userRepository,
	}
}

func (p *PushCommandHandler) Cmd() BotCmd {
	return cmd.PUSH
}

func (p *PushCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	if !b.IsBotAdmin(message.From.ID) {
		return fmt.Errorf(tip.NotAdminTip)
	}

	text := tip.DonateTip

	if utils.IsNotEmpty(message.CommandArguments()) {
		text = message.CommandArguments()
	}

	userIDs, err := p.userRepository.GetAllUserID()
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		go func(userID string, text string) {
			if utils.IsEmpty(userID) {
				return
			}
			uid, _ := utils.StringToInt64(userID)
			msg := tgbotapi.NewMessage(uid, text)
			msg.ParseMode = tgbotapi.ModeMarkdown
			b.SafeSend(msg)
		}(userID, text)
	}

	return nil
}
