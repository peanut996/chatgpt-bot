package limiter

import (
	"chatgpt-bot/bot/telegram"
	botError "chatgpt-bot/constant/error"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type JoinLimiter struct{}

func (j *JoinLimiter) Allow(bot telegram.TelegramBot, message tgbotapi.Message) (bool, string) {
	if message.Chat.IsPrivate() && !bot.Config().PrivateChatLimiter {
		return true, ""
	}

	if !message.Chat.IsPrivate() && !bot.Config().GroupChatLimiter {
		return true, ""
	}
	userID := message.From.ID
	// 限制用户加群
	config := bot.Config()
	groupName := config.BotConfig.TelegramGroupName
	channelName := config.BotConfig.TelegramChannelName
	groupUrl := "https://t.me/" + groupName
	channelUrl := "https://t.me/" + channelName
	ok := findMemberFromChat(bot, groupName, userID) &&
		findMemberFromChat(bot, channelName, userID)
	if !ok {
		msg := tgbotapi.NewMessage(message.Chat.ID, botError.LimitUserGroupAndChannel)
		button1 := tgbotapi.InlineKeyboardButton{
			URL:  &channelUrl,
			Text: "频道(Channel)",
		}
		button2 := tgbotapi.InlineKeyboardButton{
			URL:  &groupUrl,
			Text: "群组(Group)",
		}
		markup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{button1, button2}}}
		msg.ReplyMarkup = markup
		return false, botError.EmptyMessage
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
	member, err := b.TGBot().GetChatMember(findUserConfig)
	if err != nil || member.Status == "left" || member.Status == "kicked" {
		log.Printf("[ShouldLimitUser] memeber should be limit. id: %d", userID)
		return false
	}
	return true
}

func NewJoinMessageLimiter() *JoinLimiter {
	return &JoinLimiter{}
}
