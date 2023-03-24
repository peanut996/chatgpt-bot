package limiter

import (
	"chatgpt-bot/bot/telegram"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type UserLimiter struct {
	userRepository *repository.UserRepository
}

func (u *UserLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	userInfo, err := bot.GetUserInfo(message.From.ID)
	if err != nil {
		return false, err.Error()
	}

	user, err := u.userRepository.GetByUserID(utils.Int64ToString(message.From.ID))

	if err != nil {
		log.Printf("[UserLimiter] get user by user id failed, err: 【%s】\n", err)
		return false, botError.InternalError
	}

	if user == nil {
		err = u.userRepository.InitUser(utils.Int64ToString(message.From.ID), userInfo.String())
		if err != nil {
			log.Printf("[UserLimiter] init user failed, err: 【%s】\n", err)
			return false, botError.InternalError
		}
	} else {
		if user.UserName != userInfo.String() {
			err := u.userRepository.UpdateUserName(userInfo.UserName, user.UserID)
			if err != nil {
				log.Printf("[UserLimiter] update name failed, err: 【%s】\n", err)
				return false, botError.InternalError
			}
		}
	}

	return true, ""

}

func (u *UserLimiter) CallBack(telegram.TelegramBot, tgbotapi.Message, bool) {
}

func NewUserLimiter(userRepository *repository.UserRepository) *UserLimiter {
	return &UserLimiter{
		userRepository: userRepository,
	}
}
