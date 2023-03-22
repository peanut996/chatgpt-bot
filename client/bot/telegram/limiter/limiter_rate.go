package limiter

import (
	"chatgpt-bot/bot/telegram"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/middleware"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

type RateLimiter struct {
	limiter *middleware.Limiter
}

func (r *RateLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	if !bot.Config().RateLimiterConfig.Enable {
		return true, ""
	}
	if !r.limiter.Allow(strconv.FormatInt(message.From.ID, 10)) {
		log.Printf("[RateLimiter] user %d is chatting with me, ignore message %s", message.From.ID, message.Text)
		text := fmt.Sprintf(botError.RateLimitMessageTemplate,
			r.limiter.GetCapacity(), r.limiter.GetDuration()/60,
			r.limiter.GetDuration()/60, r.limiter.GetCapacity())
		return false, text
	}
	return true, ""
}

func (r *RateLimiter) CallBack(telegram.TelegramBot, tgbotapi.Message, bool) {
}

func NewRateLimiter(capacity int64, duration int64) *RateLimiter {
	return &RateLimiter{
		limiter: middleware.NewLimiter(capacity, duration),
	}
}
