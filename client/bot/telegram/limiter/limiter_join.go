package limiter

import (
	"chatgpt-bot/bot/telegram"
	botError "chatgpt-bot/constant/error"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

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

func NewJoinMessageLimiter() *JoinLimiter {
	return &JoinLimiter{}
}
