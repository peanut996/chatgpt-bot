package model

import (
	"chatgpt-bot/constant/cmd"
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
	return fmt.Sprintf("[ChatTask] [ chat: %d, from: %d, message id: %d, question: %s, answer: %s,]",
		c.Chat, c.From, c.MessageID, c.Question, c.Answer)
}

func (c *ChatTask) GetFormattedQuestion() string {
	return fmt.Sprintf("❓ from %s\n%s", c.User.String(), c.Question)
}

func (c *ChatTask) GetFormattedAnswer() string {
	return fmt.Sprintf("🅰️ to %s\n%s", c.User.String(), c.Answer)
}

func NewChatTask(message tgbotapi.Message) *ChatTask {
	task := &ChatTask{
		Question:   message.Text,
		Chat:       message.Chat.ID,
		From:       message.From.ID,
		MessageID:  message.MessageID,
		UUID:       uuid.New().String(),
		rawMessage: message,
	}
	if message.IsCommand() && message.Command() == cmd.GPT {
		task.Question = message.CommandArguments()
	}
	return task
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
