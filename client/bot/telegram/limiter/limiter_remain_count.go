package limiter

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/config"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type RemainCountMessageLimiter struct {
	userRepository *repository.UserRepository
}

func (l *RemainCountMessageLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	userID := message.From.ID
	userIDString := utils.Int64ToString(userID)

	user, err := l.userRepository.GetByUserID(userIDString)
	if err != nil {
		return false, botError.InternalError
	}
	if user.Donated() {
		return true, ""
	}

	ok := user.RemainCount > 0
	// 限制用户使用次数
	if !ok {
		user, err := l.userRepository.GetByUserID(utils.Int64ToString(userID))
		code := user.InviteCode
		link := bot.GetBotInviteLink(code)
		if err != nil {
			return false, botError.InternalError
		}
		return false, fmt.Sprintf(botError.LimitUserCountTemplate, config.CountWhenInviteOtherUser, link, config.CountWhenInviteOtherUser, link)
	}

	return true, ""
}

func (l *RemainCountMessageLimiter) CallBack(_ telegram.TelegramBot, message tgbotapi.Message, success bool) {
	if success {
		err := l.userRepository.DecreaseCount(utils.Int64ToString(message.From.ID))
		if err != nil {
			log.Println("[CallBack] decrease user count error")
			return
		}
	}

}

func NewRemainCountMessageLimiter(userRepository *repository.UserRepository) *RemainCountMessageLimiter {
	return &RemainCountMessageLimiter{
		userRepository: userRepository,
	}
}
