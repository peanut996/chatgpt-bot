package limiter

import (
	"chatgpt-bot/bot/telegram"
	botError "chatgpt-bot/constant/error"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type SingletonMessageLimiter struct {
	session *sync.Map
}

func (l *SingletonMessageLimiter) Allow(_ telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	_, ok := l.session.Load(message.From.ID)
	if ok {
		return false, botError.OnlyOneChatAtATime
	}
	defer l.session.Store(message.From.ID, true)
	return true, ""
}

func (l *SingletonMessageLimiter) CallBack(_ telegram.TelegramBot, message tgbotapi.Message, _ bool) {
	l.session.Delete(message.From.ID)
}

func NewSingletonMessageLimiter() *SingletonMessageLimiter {
	return &SingletonMessageLimiter{
		session: &sync.Map{},
	}
}
