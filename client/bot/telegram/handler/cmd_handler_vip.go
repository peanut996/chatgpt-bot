package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type VIPCommandHandler struct {
	userRepository *repository.UserRepository
}

func (v *VIPCommandHandler) Cmd() BotCmd {
	return cmd.VIP
}

func (v *VIPCommandHandler) Run(t telegram.TelegramBot, message tgbotapi.Message) error {
	if !t.IsBotAdmin(message.From.ID) {
		return fmt.Errorf(tip.NotAdminTip)
	}
	args := message.CommandArguments()
	if args == "" {
		return fmt.Errorf(botError.MissingRequiredConfig + "id")
	}
	err := v.userRepository.UpdateUserToVIP(message.From.ID)
	if err != nil {
		return err
	}
	return fmt.Errorf("success")
}

func NewVIPCommandHandler(userRepository *repository.UserRepository) *VIPCommandHandler {
	return &VIPCommandHandler{userRepository: userRepository}
}
