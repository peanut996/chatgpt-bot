package limiter

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/bot/telegram/service"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/middleware"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Limiter interface {
	Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string)
	CallBack(bot telegram.TelegramBot, message tgbotapi.Message, success bool)
}

type UserLimiter struct {
	userRepository *repository.UserRepository
}

func (u *UserLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	user, err := u.userRepository.GetByUserID(utils.Int64ToString(message.From.ID))

	if err != nil {
		log.Printf("[UserLimiter] get user by user id failed, err: 【%s】\n", err)
		return false, botError.InternalError
	}

	if user == nil {
		userInfo, err := bot.GetUserInfo(message.From.ID)
		if err != nil {
			return false, err.Error()
		}
		err = u.userRepository.InitUser(utils.Int64ToString(message.From.ID), userInfo.String())
		if err != nil {
			log.Printf("[UserLimiter] init user failed, err: 【%s】\n", err)
			return false, botError.InternalError
		}
	}

	return true, ""

}

func (u *UserLimiter) CallBack(telegram.TelegramBot, tgbotapi.Message, bool) {
}

type CommonMessageLimiter struct {
}

func (l *CommonMessageLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
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
		!(message.ReplyToMessage.From.ID == bot.SelfID()) {
		// 不是回复机器人的不用处理
		return false, botError.EmptyMessage
	}

	isPrivate := message.Chat.IsPrivate()
	// 私聊或者是回复机器人的消息才处理或者是机器人的命令
	ok := isPrivate || service.IsGPTMessage(message) ||
		(message.ReplyToMessage != nil &&
			message.ReplyToMessage.From.ID == bot.SelfID())

	return ok, ""
}

func (l *CommonMessageLimiter) CallBack(b telegram.TelegramBot, m tgbotapi.Message, success bool) {
	shouldSendTip := func() bool {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := r.Intn(100)
		return n == 8
	}
	if success && m.Chat.IsPrivate() && shouldSendTip() {
		go func() {
			time.Sleep(time.Second * 30)
			b.SafeSendMsg(m.Chat.ID, tip.DonateTip)
		}()
	}
}

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

type JoinLimiter struct{}

func (j *JoinLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	if !message.Chat.IsPrivate() {
		return true, ""
	}
	userID := message.From.ID
	// 限制用户加群
	config := bot.Config()
	groupName := config.BotConfig.TelegramGroupName
	channelName := config.BotConfig.TelegramChannelName
	if config.BotConfig.ShouldLimitPrivate {
		ok := findMemberFromChat(bot, groupName, userID) &&
			findMemberFromChat(bot, channelName, userID)
		if !ok {
			return false, fmt.Sprintf(botError.LimitUserGroupAndChannelTemplate,
				channelName, groupName, channelName, groupName)
		}
	}
	return true, ""
}

func (j *JoinLimiter) CallBack(telegram.TelegramBot, tgbotapi.Message, bool) {
}

func findMemberFromChat(b telegram.TelegramBot, chatName string, userID int64) bool {
	findUserConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			SuperGroupUsername: chatName,
			UserID:             userID,
		},
	}
	member, err := b.GetAPIBot().GetChatMember(findUserConfig)
	if err != nil || member.Status == "left" || member.Status == "kicked" {
		log.Printf("[ShouldLimitUser] memeber should be limit. id: %d", userID)
		return false
	}
	return true
}

func (l *RemainCountMessageLimiter) CallBack(_ telegram.TelegramBot, message tgbotapi.Message, success bool) {
	if !service.IsGPTMessage(message) {
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

func NewCommonMessageLimiter() *CommonMessageLimiter {
	return &CommonMessageLimiter{}
}

func NewRemainCountMessageLimiter(userRepository *repository.UserRepository) *RemainCountMessageLimiter {
	return &RemainCountMessageLimiter{
		userRepository: userRepository,
	}
}

func NewSingletonMessageLimiter() *SingletonMessageLimiter {
	return &SingletonMessageLimiter{
		session: &sync.Map{},
	}
}

func NewJoinMessageLimiter() *JoinLimiter {
	return &JoinLimiter{}
}

func NewUserLimiter(userRepository *repository.UserRepository) *UserLimiter {
	return &UserLimiter{
		userRepository: userRepository,
	}
}
