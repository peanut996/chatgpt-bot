package limiter

import (
	"chatgpt-bot/bot/telegram"
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

	// 查看用户是否存在 不存在就初始化
	user, err := l.userRepository.GetByUserID(userIDString)
	if err != nil {
		return false, botError.InternalError
	}
	if user == nil {
		// 初始化用户
		userName := ""
		tgUser, err := bot.GetUserInfo(message.From.ID)
		if err == nil {
			userName = tgUser.String()
		}
		err = l.userRepository.InitUser(userIDString, userName)
		if err != nil {
			log.Println("RemainCountMessageLimiter] init user error", err)
			return false, botError.InternalError
		}
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
		return false, fmt.Sprintf(botError.LimitUserCountTemplate, link, link)
	}

	return true, ""
}

func (l *RemainCountMessageLimiter) CallBack(_ telegram.TelegramBot, message tgbotapi.Message, success bool) {
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

func NewRemainCountMessageLimiter(userRepository *repository.UserRepository) *RemainCountMessageLimiter {
	return &RemainCountMessageLimiter{
		userRepository: userRepository,
	}
}
