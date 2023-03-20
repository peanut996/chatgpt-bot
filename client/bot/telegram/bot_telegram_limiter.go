package telegram

import (
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/middleware"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageLimiter interface {
	Allow(bot *Bot, message tgbotapi.Message) (bool, string)
	CallBack(bot *Bot, message tgbotapi.Message, success bool)
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
		return false, botError.EmptyMessage
	}
	if strings.Trim(message.Text, " ") == "" || (message.IsCommand() &&
		strings.Trim(message.CommandArguments(), " ") == "") {
		// 空消息不用处理
		return false, botError.EmptyMessage
	}

	if message.ReplyToMessage != nil &&
		!(message.ReplyToMessage.From.ID == bot.tgBot.Self.ID) {
		// 不是回复机器人的不用处理
		return false, botError.EmptyMessage
	}

	isPrivate := message.Chat.IsPrivate()
	// 私聊或者是回复机器人的消息才处理
	ok := isPrivate ||
		(message.ReplyToMessage != nil &&
			message.ReplyToMessage.From.ID == bot.tgBot.Self.ID)

	return ok, ""
}

func (l *CommonMessageLimiter) CallBack(*Bot, tgbotapi.Message, bool) {
}

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
		return false, botError.OnlyOneChatAtATime
	}
	defer l.session.Store(message.From.ID, true)
	return true, ""
}

func (l *SingletonMessageLimiter) CallBack(_ *Bot, message tgbotapi.Message, _ bool) {
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
	if !message.Chat.IsPrivate() {
		return true, ""
	}
	userID := message.From.ID
	userIDString := utils.Int64ToString(userID)
	// 限制用户加群
	if bot.limitPrivate {
		ok := findMemberFromChat(bot, bot.groupName, userID) &&
			findMemberFromChat(bot, bot.channelName, userID)
		if !ok {
			return false, fmt.Sprintf(botError.LimitUserGroupAndChannelTemplate,
				bot.channelName, bot.groupName, bot.channelName, bot.groupName)
		}
	}

	// 查看用户是否存在 不存在就初始化
	user, err := l.userRepository.GetByUserID(userIDString)
	if err != nil {
		return false, botError.InternalError
	}
	if user == nil {
		// 初始化用户
		userName := ""
		tgUser, err := bot.getUserInfo(message.From.ID)
		if err == nil {
			userName = tgUser.String()
		}
		err = l.userRepository.InitUser(userIDString, userName)
		if err != nil {
			log.Println("PrivateMessageLimiter] init user error", err)
			return false, botError.InternalError
		}
		return true, ""
	}

	ok := user.RemainCount > 0
	// 限制用户使用次数
	if !ok {
		user, err := l.userRepository.GetByUserID(utils.Int64ToString(userID))
		code := user.InviteCode
		link := bot.getBotInviteLink(code)
		if err != nil {
			return false, botError.InternalError
		}
		return false, fmt.Sprintf(botError.LimitUserCountTemplate, link, link)
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

func (l *PrivateMessageLimiter) CallBack(_ *Bot, message tgbotapi.Message, success bool) {
	if !IsGPT4Message(message) {
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
		text := fmt.Sprintf(botError.RateLimitMessageTemplate,
			r.limiter.GetCapacity(), r.limiter.GetDuration()/60,
			r.limiter.GetDuration()/60, r.limiter.GetCapacity())
		return false, text
	}
	return true, ""
}

func (r *RateLimiter) CallBack(*Bot, tgbotapi.Message, bool) {
}
