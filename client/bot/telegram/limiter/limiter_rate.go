package limiter

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

type RateLimiter struct {
	limiter *middleware.Limiter
}

func (r *RateLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	if !bot.Config().Downgrade {
		return true, ""
	}
	allow, err := r.limiter.Allow(strconv.FormatInt(message.From.ID, 10))
	if !allow {
		log.Printf("[RateLimiter] user %d is chatting with me, ignore message %s", message.From.ID, message.Text)
		return false, err.Error()
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
