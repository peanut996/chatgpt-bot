package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/config"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AccessCommandHandler struct {
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
	salt                       string
}

func (a *AccessCommandHandler) Cmd() BotCmd {
	return cmd.ACCESS
}

func (a *AccessCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	if !message.Chat.IsPrivate() {
		b.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, botError.OnlyAllowInPrivate)
		return nil
	}
	user, err := a.userRepository.GetByUserID(utils.Int64ToString(message.From.ID))
	if err != nil {
		b.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, botError.InternalError)
		return nil
	}
	if user.Donated() {
		accessCode := utils.GetAccessCode(user.InviteCode, a.salt)
		b.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, fmt.Sprintf(tip.AccessCodeTipTemplate, accessCode, accessCode))
		return nil
	}

	count, err := a.userInviteRecordRepository.CountByUserID(user.UserID)
	if err != nil {
		b.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, botError.InternalError)
		return nil
	}
	if count >= int64(config.AllowByInviteCount) {
		accessCode := utils.GetAccessCode(user.InviteCode, a.salt)
		b.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, fmt.Sprintf(tip.AccessCodeTipTemplate, accessCode, accessCode))
		return nil
	}
	cannotGetAccessTip := fmt.Sprintf(botError.CannotGetAccessCodeTemplate, config.AllowByInviteCount, config.AllowByInviteCount)
	b.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, cannotGetAccessTip)
	return nil
}
func NewAccessCommandHandler(userRepository *repository.UserRepository,
	userInviteRecordRepository *repository.UserInviteRecordRepository,
	salt string) *AccessCommandHandler {
	return &AccessCommandHandler{
		salt:                       salt,
		userRepository:             userRepository,
		userInviteRecordRepository: userInviteRecordRepository,
	}
}
