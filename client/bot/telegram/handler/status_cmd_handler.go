package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StatusCommandHandler struct {
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
}

func (s *StatusCommandHandler) Cmd() BotCmd {
	return cmd.STATUS
}

func (s *StatusCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	if !b.IsBotAdmin(message.From.ID) {
		return nil
	}
	userCount, err := s.userRepository.Count()
	if err != nil {
		return err
	}

	inviteRecordCount, err := s.userInviteRecordRepository.Count()
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(tip.StatusTipTemplate, userCount, inviteRecordCount))
	b.SafeSend(msg)
	return nil
}

func NewStatusCommandHandler(userRepository *repository.UserRepository, userInviteRecordRepository *repository.UserInviteRecordRepository) *StatusCommandHandler {
	return &StatusCommandHandler{
		userRepository:             userRepository,
		userInviteRecordRepository: userInviteRecordRepository,
	}
}
