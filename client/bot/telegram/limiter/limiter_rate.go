package limiter

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/config"
	"chatgpt-bot/middleware"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type RateLimiter struct {
	isGPT4Limiter              bool
	limiter                    *middleware.Limiter
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
}

func (r *RateLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	if !bot.Config().EnableRateLimiter {
		return true, ""
	}
	userIDString := utils.Int64ToString(message.From.ID)
	user, err := r.userRepository.GetByUserID(userIDString)
	if err != nil {
		log.Printf("[RateLimiter] error when get user %d: %s", message.From.ID, err.Error())
		return false, err.Error()
	}
	if user.Donated() {
		return true, ""
	}

	if !r.isGPT4Limiter {
		count, err := r.userInviteRecordRepository.CountByUserID(user.UserID)
		if err != nil {
			log.Printf("[RateLimiter] error when get user %d: %s", message.From.ID, err.Error())
			return false, err.Error()
		}

		if count > int64(config.AllowByInviteCount) {
			return true, ""
		}
	}
	allow, err := r.limiter.Allow(userIDString)
	if !allow {
		log.Printf("[RateLimiter] user %d is chatting with me, ignore message %s", message.From.ID, message.Text)
		return false, err.Error()
	}
	return true, ""
}

func (r *RateLimiter) CallBack(telegram.TelegramBot, tgbotapi.Message, bool) {
}

func NewRateLimiter(capacity int64, duration int64,
	isGPT4Limiter bool,
	userRepository *repository.UserRepository,
	recordRepository *repository.UserInviteRecordRepository) *RateLimiter {
	return &RateLimiter{
		limiter:                    middleware.NewLimiter(capacity, duration),
		userRepository:             userRepository,
		userInviteRecordRepository: recordRepository,
		isGPT4Limiter:              isGPT4Limiter,
	}
}
