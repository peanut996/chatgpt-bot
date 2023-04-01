package limiter

import (
	"chatgpt-bot/bot/telegram"
	botError "chatgpt-bot/constant/error"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		!(message.ReplyToMessage.From.ID == bot.SelfID()) && !IsGPTMessage(message) {
		// 不是回复机器人的不用处理
		return false, botError.EmptyMessage
	}

	isPrivate := message.Chat.IsPrivate()
	// 私聊或者是回复机器人的消息才处理或者是机器人的命令
	ok := isPrivate || IsGPTMessage(message) ||
		(message.ReplyToMessage != nil &&
			message.ReplyToMessage.From.ID == bot.SelfID())

	return ok, ""
}

func (l *CommonMessageLimiter) CallBack(telegram.TelegramBot, tgbotapi.Message, bool) {

}

func NewCommonMessageLimiter() *CommonMessageLimiter {
	return &CommonMessageLimiter{}
}
