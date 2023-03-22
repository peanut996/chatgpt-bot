package limiter

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/config"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type InviteCountMessageLimiter struct {
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
}

func (l *InviteCountMessageLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	userID := message.From.ID
	userIDString := utils.Int64ToString(userID)

	count, err := l.userInviteRecordRepository.CountByUserID(userIDString)
	if err != nil {
		return false, botError.InternalError
	}
	ok := count >= int64(config.AllowGPT4Count)
	// 限制用户使用次数
	if !ok {
		user, err := l.userRepository.GetByUserID(utils.Int64ToString(userID))
		code := user.InviteCode
		link := bot.GetBotInviteLink(code)
		if err != nil {
			return false, botError.InternalError
		}
		return false, fmt.Sprintf(tip.InviteTipTemplate, config.AllowGPT4Count, link,
			config.AllowGPT4Count, link)
	}

	return true, ""
}

func (l *InviteCountMessageLimiter) CallBack(_ telegram.TelegramBot, message tgbotapi.Message, success bool) {
	if !IsGPTMessage(message) {
		return
	}
	if success {
		err := l.userRepository.DecreaseCount(utils.Int64ToString(message.From.ID))
		if err != nil {
			log.Println("[CallBack] decrease user count error")
			return
		}
	}

}

func NewInviteCountLimiter(userRepository *repository.UserRepository,
	userInviteRecordRepository *repository.UserInviteRecordRepository) *InviteCountMessageLimiter {
	return &InviteCountMessageLimiter{
		userRepository:             userRepository,
		userInviteRecordRepository: userInviteRecordRepository,
	}
}
