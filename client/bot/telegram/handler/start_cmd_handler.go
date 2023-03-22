package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/model/persist"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type StartCommandHandler struct {
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
}

func (h *StartCommandHandler) Cmd() BotCmd {
	return cmd.START
}

func (h *StartCommandHandler) Run(bot telegram.TelegramBot, message tgbotapi.Message) error {
	log.Println(fmt.Printf("get args: [%h]", message.CommandArguments()))
	args := message.CommandArguments()
	if matchInviteCode(args) {
		err := h.handleInvitation(args, utils.Int64ToString(message.From.ID), bot)
		if err != nil {
			log.Printf("[StartCommandHandler] handle invitation failed, err: 【%h】", err)
		}
	}
	bot.SafeSendMsg(message.Chat.ID, tip.BotStartTip)
	return nil
}

func (h *StartCommandHandler) handleInvitation(inviteCode string, inviteUserID string, b telegram.TelegramBot) error {
	user, err := h.userRepository.GetUserByInviteCode(inviteCode)
	if err != nil {
		log.Printf("[handleInvitation] find user by invite code failed, err: 【%h】", err)
		return err
	}
	if user == nil {
		log.Printf("[handleInvitation] find user by invite code failed, user is nil")
		return errors.New("no such user by invite code: " + inviteCode)
	}
	if user.UserID == inviteUserID {
		log.Printf("[handleInvitation] user can not invite himself")
		return fmt.Errorf("[handleInvitation] user can not invite himself, user id: [%h]", inviteUserID)
	}
	record, err := h.userInviteRecordRepository.GetByInviteUserID(inviteUserID)
	if err != nil {
		log.Printf("[handleInvitation] find user by invite user id failed, err: 【%h】", err)
		return err
	}
	if record != nil {
		log.Printf("[handleInvitation]  user has been invited by other user: " + record.UserID)
		return nil
	}
	inviteRecord := persist.NewUserInviteRecord(user.UserID, inviteUserID)
	err = h.userInviteRecordRepository.Insert(inviteRecord)
	if err != nil {
		return err
	}
	err = h.userRepository.AddCountWhenInviteOther(user.UserID)
	if err != nil {
		return err
	}
	originUserID, _ := utils.StringToInt64(user.UserID)
	b.SafeSendMsg(originUserID, tip.InviteSuccessTip)
	return nil
}

func matchInviteCode(code string) bool {
	return utils.IsNotEmpty(code) && len(code) == 10 && utils.IsMatchString(`^[a-zA-Z]{10}$`, code)
}

func NewStartCommandHandler(userRepository *repository.UserRepository, userInviteRecordRepository *repository.UserInviteRecordRepository) *StartCommandHandler {
	return &StartCommandHandler{
		userRepository:             userRepository,
		userInviteRecordRepository: userInviteRecordRepository,
	}
}
