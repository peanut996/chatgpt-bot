package model

import "fmt"

type ChatTask struct {
	Question  string
	Answer    string
	Chat      int64
	From      int64
	MessageID int
}

func (c *ChatTask) String() string {
	return fmt.Sprintf("[ChatTask] [question: %s, answer: %s, chat: %d, from: %d, message id: %d]", c.Question, c.Answer, c.Chat, c.From, c.MessageID)
}

func (c *ChatTask) GetFormattedQuestion() string {
	return fmt.Sprintf("❓ %s", c.Question)
}

func (c *ChatTask) GetFormattedAnswer() string {
	return fmt.Sprintf("🅰️ %s", c.Answer)
}

func NewChatTask(question string, chat, from int64, msgID int) *ChatTask {
	return &ChatTask{
		Question:  question,
		Chat:      chat,
		From:      from,
		MessageID: msgID,
	}
}
