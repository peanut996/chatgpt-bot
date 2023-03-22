package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/config"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type QueryCommandHandler struct {
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
}

func (q *QueryCommandHandler) Cmd() BotCmd {
	return cmd.QUERY
}

func (q *QueryCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	userID := utils.Int64ToString(message.From.ID)
	user, err := q.userRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("[QueryCommandHandler] get user by user id failed, err: 【%s】\n", err)
		return err
	}
	if user == nil {
		userInfo, err := b.GetUserInfo(message.From.ID)
		if err != nil {
			return err
		}
		err = q.userRepository.InitUser(userID, userInfo.String())
		if err != nil {
			log.Printf("[QueryCommandHandler] init user failed, err: 【%s】\n", err)
			return err
		}
		user, err = q.userRepository.GetByUserID(userID)
		if err != nil {
			log.Printf("[QueryCommandHandler] get user by user id failed, err: 【%s】\n", err)
			return err
		}
	}
	inviteCount, err := q.userInviteRecordRepository.CountByUserID(userID)
	if err != nil {
		log.Printf("[QueryCommandHandler] get user invite count by user id failed, err: 【%s】\n", err)
		return err
	}

	text := fmt.Sprintf(tip.QueryUserInfoTemplate,
		userID, user.RemainCount, inviteCount, b.GetBotInviteLink(user.InviteCode),
		config.AllowGPT4Count, config.AllowGPT4Count)
	b.SafeReplyMsg(message.Chat.ID, message.MessageID, text)
	return nil
}

func NewQueryCommandHandler(userRepository *repository.UserRepository, userInviteRecordRepository *repository.UserInviteRecordRepository) *QueryCommandHandler {
	return &QueryCommandHandler{
		userRepository,
		userInviteRecordRepository,
	}
}
