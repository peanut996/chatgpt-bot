package model

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type ChatTask struct {
	Question  string
	Answer    string
	Chat      int64
	From      int64
	MessageID int
	UUID      string

	User          *User
	rawMessage    tgbotapi.Message
	IsGPT4Message bool
}

func (c *ChatTask) String() string {
	return fmt.Sprintf("[ChatTask] [question: %s, answer: %s, chat: %d, from: %d, message id: %d]", c.Question, c.Answer, c.Chat, c.From, c.MessageID)
}

func (c *ChatTask) GetFormattedQuestion() string {
	return fmt.Sprintf("❓ from %s\n%s", c.User.String(), c.Question)
}

func (c *ChatTask) GetFormattedAnswer() string {
	return fmt.Sprintf("🅰️ to %s\n%s", c.User.String(), c.Answer)
}

func NewChatTask(message tgbotapi.Message) *ChatTask {
	return &ChatTask{
		Question:   message.Text,
		Chat:       message.Chat.ID,
		From:       message.From.ID,
		MessageID:  message.MessageID,
		UUID:       uuid.New().String(),
		rawMessage: message,
	}
}

func NewGPT4ChatTask(message tgbotapi.Message) *ChatTask {
	chatTask := NewChatTask(message)
	chatTask.Question = message.CommandArguments()
	chatTask.IsGPT4Message = true
	return chatTask
}

func (c *ChatTask) GetRawMessage() tgbotapi.Message {
	return c.rawMessage
}
