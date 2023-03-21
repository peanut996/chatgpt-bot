package telegram

import (
	"chatgpt-bot/constant/cmd"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/model"
	"chatgpt-bot/utils"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsGPTMessage(message tgbotapi.Message) bool {
	return message.IsCommand() && (message.Command() == cmd.GPT4 || message.Command() == cmd.GPT)
}
func IsGPT3Message(message tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == cmd.GPT
}

func (b *Bot) isBotAdmin(from int64) bool {
	if b.admin == 0 {
		return false
	}
	return b.admin == from
}

func (b *Bot) getBotInviteLink(code string) string {
	return fmt.Sprintf("https://t.me/%s?start=%s", b.tgBot.Self.UserName, code)
}

func (b *Bot) getUserInfo(userID int64) (*model.User, error) {
	user, err := b.tgBot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: userID,
		}})
	if err != nil {
		log.Println("[GetUserInfo] get user info failed, err: ", err)
		return nil, err
	}
	log.Println("[GetUserInfo] get user info success: ", utils.ToIndentJson(user))
	return model.NewUser(user.UserName, user.FirstName, user.LastName), nil
}

func (b *Bot) sendTyping(chatID int64) {
	msg := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
	_, _ = b.tgBot.Send(msg)
}

func (b *Bot) sendErrorMessage(chatID int64, msgID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = msgID
	_, err := b.tgBot.Send(msg)
	if err != nil {
		log.Printf("[SendErrorMessage] send message failed, err: 【%s】, msg: 【%+v】", err, msg)
		msg.Text = botError.SendBackMsgFailed
		_, _ = b.tgBot.Send(msg)
		return
	}
}

func (b *Bot) safeSend(msg tgbotapi.MessageConfig) {
	if msg.Text == "" {
		return
	}
	if len(msg.Text) < 4096 {
		b.sendMessageSilently(msg)
		return
	}
	b.sendLargeMessage(msg)
}

func (b *Bot) sendLargeMessage(msg tgbotapi.MessageConfig) {
	msgs := utils.SplitMessageByMaxSize(msg.Text, 4000)
	for _, m := range msgs {
		msg.Text = m
		b.sendMessageSilently(msg)
	}
}

func (b *Bot) sendMessageSilently(msg tgbotapi.MessageConfig) {
	if msg.Text == "" {
		return
	}
	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err := b.tgBot.Send(msg)
	if err != nil {
		log.Printf("[SendMessageSilently] send message failed, err: 【%s】, msg: 【%+v】", err, msg)
		msg.ParseMode = ""
		_, _ = b.tgBot.Send(msg)
	}
}

func (b *Bot) sendFromChatTask(task model.ChatTask) {
	msg := tgbotapi.NewMessage(task.Chat, task.Question)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.Text = task.Answer
	msg.ReplyToMessageID = task.MessageID
	msgs := utils.SplitMessageByMaxSize(task.Answer, 4000)
	for _, m := range msgs {
		msg.Text = m
		b.safeSend(msg)
	}
}

func (b *Bot) safeSendMsg(chatID int64, text string) {
	b.safeSend(tgbotapi.NewMessage(chatID, text))
}

func (b *Bot) safeReplyMsg(chatID int64, messageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	b.safeSend(msg)
}

func (b *Bot) logToChannel(log string) {
	go func(s string) {
		msg := tgbotapi.NewMessage(b.logChannelID, s)
		msg.ParseMode = tgbotapi.ModeMarkdown
		b.safeSend(msg)
	}(log)
}

func (b *Bot) sendQueueToast(chatID int64, messageID int) {
	queue := len(b.chatTaskChannel)
	if queue < 3 {
		return
	}
	text := fmt.Sprintf(tip.QueueTipTemplate, queue, queue)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	msg.ParseMode = tgbotapi.ModeMarkdown
	b.safeSend(msg)

}
