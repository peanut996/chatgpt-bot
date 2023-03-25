package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type InviteCommandHandler struct {
	userRepository *repository.UserRepository
}

func (i *InviteCommandHandler) Cmd() BotCmd {
	return cmd.INVITE
}

func (i *InviteCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	userID := utils.Int64ToString(message.From.ID)
	user, err := i.userRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("[InviteCommandHandler] find user by user id failed, err: 【%s】", err)
		return err
	}
	if user != nil {
		link := b.GetBotInviteLink(user.InviteCode)
		b.SafeSendMsgWithoutPreview(message.Chat.ID, fmt.Sprintf(tip.InviteLinkTemplate, link, link))
		return nil
	} else {
		userName := ""
		tgUser, err := b.GetUserInfo(message.From.ID)
		if err == nil {
			userName = tgUser.String()
		}
		err = i.userRepository.InitUser(userID, userName)
		if err != nil {
			log.Printf("[InviteCommandHandler] init user failed, err: 【%s】", err)
			return err
		}
		user, _ := i.userRepository.GetByUserID(userID)
		link := b.GetBotInviteLink(user.InviteCode)
		b.SafeSendMsgWithoutPreview(message.Chat.ID, fmt.Sprintf(tip.InviteLinkTemplate, link, link))
	}
	return nil
}

func NewInviteCommandHandler(userRepository *repository.UserRepository) *InviteCommandHandler {
	return &InviteCommandHandler{
		userRepository: userRepository,
	}
}
