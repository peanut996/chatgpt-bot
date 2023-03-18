package telegram

import (
	"chatgpt-bot/constant"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type MessageLimiter interface {
	Allow(bot *Bot, message tgbotapi.Message) (bool, string)
}

type CommonMessageLimiter struct {
}

func NewCommonMessageLimiter() *CommonMessageLimiter {
	return &CommonMessageLimiter{}
}

func (l *CommonMessageLimiter) Allow(bot *Bot, message tgbotapi.Message) (bool, string) {
	isPrivate := message.Chat.IsPrivate()
	ok := isPrivate ||
		(message.ReplyToMessage != nil &&
			message.ReplyToMessage.From.ID == bot.tgBot.Self.ID)
	return ok, ""
}

type SingleMessageLimiter struct {
}

func NewSingleMessageLimiter() *SingleMessageLimiter {
	return &SingleMessageLimiter{}
}

func (l *SingleMessageLimiter) Allow(*Bot, tgbotapi.Message) (bool, string) {
	return true, ""
}

type PrivateMessageLimiter struct {
	userRepository repository.UserRepository
}

func NewPrivateMessageLimiter(userRepository repository.UserRepository) *PrivateMessageLimiter {
	return &PrivateMessageLimiter{
		userRepository: userRepository,
	}
}

func (l *PrivateMessageLimiter) Allow(bot *Bot, message tgbotapi.Message) (bool, string) {
	userID := message.From.ID
	ok, err := l.userRepository.IsAvaliable(utils.ConvertUserID(userID))
	if err != nil {
		return false, err.Error()
	}
	if !ok {
		link, err := l.userRepository.GetUserInviteLink(utils.ConvertUserID(userID))
		if err != nil {
			return false, constant.InternalError
		}
		return false, fmt.Sprintf(constant.LimitUserMessageTemplate, link, link)
	}
	return true, ""
}

type RateLimiter struct {
	maxMessages int
	maxInterval int
	lastMessage time.Time
}

func (r *RateLimiter) Allow(*Bot, tgbotapi.Message) (bool, string) {
	return true, ""
}
