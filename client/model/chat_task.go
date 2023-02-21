package model

type ChatTask struct {
	Question  string
	Answer    string
	Chat      int64
	From      int64
	MessageID int
}

func NewChatTask(question string, chat, from int64, msgID int) *ChatTask {
	return &ChatTask{
		Question:  question,
		Chat:      chat,
		From:      from,
		MessageID: int(msgID),
	}
}
