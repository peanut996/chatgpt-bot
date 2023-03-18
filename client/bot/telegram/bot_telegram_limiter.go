package telegram

import (
	"chatgpt-bot/constant"
	"chatgpt-bot/middleware"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"sync"
)

type MessageLimiter interface {
	Allow(bot *Bot, message tgbotapi.Message) (bool, string)
	CallBack(bot *Bot, message tgbotapi.Message)
}

type CommonMessageLimiter struct {
}

func NewCommonMessageLimiter() *CommonMessageLimiter {
	return &CommonMessageLimiter{}
}

func (l *CommonMessageLimiter) Allow(bot *Bot, message tgbotapi.Message) (bool, string) {
	if message.NewChatMembers != nil ||
		message.LeftChatMember != nil {
		// 新成员加入或者成员离开不用处理
		return false, constant.EmptyMessage
	}
	if strings.Trim(message.Text, " ") == "" {
		// 空消息不用处理
		return false, constant.EmptyMessage
	}

	if message.ReplyToMessage != nil &&
		!(message.ReplyToMessage.From.ID == bot.tgBot.Self.ID) {
		// 不是回复机器人的不用处理
		return false, constant.EmptyMessage
	}

	isPrivate := message.Chat.IsPrivate()
	// 私聊或者是回复机器人的消息才处理
	ok := isPrivate ||
		(message.ReplyToMessage != nil &&
			message.ReplyToMessage.From.ID == bot.tgBot.Self.ID)

	return ok, ""
}

func (l *CommonMessageLimiter) CallBack(*Bot, tgbotapi.Message) {
}

// SingletonMessageLimiter allows only one message at a time
type SingletonMessageLimiter struct {
	session *sync.Map
}

func NewSingleMessageLimiter() *SingletonMessageLimiter {
	return &SingletonMessageLimiter{
		session: &sync.Map{},
	}
}

func (l *SingletonMessageLimiter) Allow(_ *Bot, message tgbotapi.Message) (bool, string) {
	_, ok := l.session.Load(message.From.ID)
	if ok {
		return false, constant.OnlyOneChatAtATime
	}
	defer l.session.Store(message.From.ID, true)
	return true, ""
}

func (l *SingletonMessageLimiter) CallBack(_ *Bot, message tgbotapi.Message) {
	l.session.Delete(message.From.ID)
}

type PrivateMessageLimiter struct {
	userRepository *repository.UserRepository
}

func NewPrivateMessageLimiter(userRepository *repository.UserRepository) *PrivateMessageLimiter {
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
	// 限制用户使用次数
	if !ok {
		link, err := l.userRepository.GetUserInviteLink(utils.ConvertUserID(userID))
		if err != nil {
			return false, constant.InternalError
		}
		return false, fmt.Sprintf(constant.LimitUserCountTemplate, link, link)
	}

	// 限制用户加群
	if bot.limitPrivate {
		ok = findMemberFromChat(bot, bot.groupName, userID)
		if !ok {
			return false, fmt.Sprintf(constant.LimitUserGroupAndChannelTemplate,
				bot.channelName, bot.groupName, bot.channelName, bot.groupName)
		}
	}

	return true, ""
}

func findMemberFromChat(b *Bot, chatName string, userID int64) bool {
	findUserConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			SuperGroupUsername: chatName,
			UserID:             userID,
		},
	}
	member, err := b.tgBot.GetChatMember(findUserConfig)
	if err != nil || member.Status == "left" || member.Status == "kicked" {
		log.Printf("[ShouldLimitUser] memeber should be limit. id: %d", userID)
		return false
	}
	return true
}

func (l *PrivateMessageLimiter) CallBack(*Bot, tgbotapi.Message) {
}

type RateLimiter struct {
	limiter *middleware.Limiter
}

func NewRateLimiter(capacity int64, duration int64) *RateLimiter {
	return &RateLimiter{
		limiter: middleware.NewLimiter(capacity, duration),
	}
}

func (r *RateLimiter) Allow(bot *Bot, message tgbotapi.Message) (bool, string) {
	if !bot.enableLimiter {
		return true, ""
	}
	if !r.limiter.Allow(strconv.FormatInt(message.From.ID, 10)) {
		log.Printf("[RateLimiter] user %d is chatting with me, ignore message %s", message.From.ID, message.Text)
		text := fmt.Sprintf(constant.RateLimitMessageTemplate,
			r.limiter.GetCapacity(), r.limiter.GetDuration()/60,
			bot.channelName, bot.groupName,
			r.limiter.GetDuration()/60, r.limiter.GetCapacity(),
			bot.channelName, bot.groupName)
		return false, text
	}
	return true, ""
}

func (r *RateLimiter) CallBack(*Bot, tgbotapi.Message) {
}
