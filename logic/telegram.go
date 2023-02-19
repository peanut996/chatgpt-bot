package logic

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type ChatTask struct {
	Question  string
	Answer    string
	Chat      int64
	From      int64
	MessageID int
}

var (
	bot           *tgbotapi.BotAPI
	offset        int = 0
	session       *sync.Map
	TaskChannel   chan *ChatTask
	tgChannelName string
	tgGroupName   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	b, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}
	bot = b
	session = &sync.Map{}
	TaskChannel = make(chan *ChatTask, 1)
	tgChannelName = os.Getenv("TELEGRAM_CHANNEL_NAME")
	tgGroupName = os.Getenv("TELEGRAM_GROUP_NAME")

	go loopAndFinishChatTask()
}

func NewChatTask(question string, chat, from int64, msgID int) *ChatTask {
	return &ChatTask{
		Question:  question,
		Chat:      chat,
		From:      from,
		MessageID: int(msgID),
	}
}

func sendTaskToChannel(question string, chat, from int64, msgID int) {
	session.Store(from, &struct{}{})
	log.Printf("[SendTaskToChannel] with question %s, chat id: %d, from: %d", question, chat, from)
	chatTask := NewChatTask(question, chat, from, msgID)
	TaskChannel <- chatTask
	chatTask.SendTyping()
}

func (t *ChatTask) Send() {
	msg := tgbotapi.NewMessage(t.Chat, t.Question)
	msg.ParseMode = "markdown"
	msg.Text = t.Answer
	msg.ReplyToMessageID = t.MessageID
	bot.Send(msg)
}

func (t *ChatTask) GetAnswerFromChatGPT() {
	a := GetChatGPTResponseWithRetry(t.Question)
	t.Answer = a
}

func (t *ChatTask) SendTyping() {
	action := tgbotapi.NewChatAction(t.Chat, tgbotapi.ChatTyping)
	bot.Send(action)
}

func (t *ChatTask) Finish() {
	log.Printf("[Finish] start chat task with question %s, chat id: %d, from: %d", t.Question, t.Chat, t.From)
	defer session.Delete(t.From)

	t.GetAnswerFromChatGPT()
	t.Send()

	log.Printf("[Finish] end chat task with question %s, chat id: %d, from: %d", t.Question, t.Chat, t.From)

}

func FetchUpdates() {
	config := tgbotapi.NewUpdate(offset)
	config.Timeout = 60

	botChannel := bot.GetUpdatesChan(config)
	for {
		select {
		case update, ok := <-botChannel:
			if !ok {
				botChannel = bot.GetUpdatesChan(config)
				log.Println("[FetchUpdates] channel closed, fetch again")
				continue
			}
			go handleUpdate(update)
		case <-time.After(30 * time.Second):
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[BotUpdate] update id:[%d] from [%s] : %s", update.UpdateID, update.Message.From.String(), update.Message.Text)

	if update.Message.IsCommand() {
		msg := handleCommandMsg(update)
		bot.Send(msg)
	} else {
		msg, hasSentChatTask := handleUserMessage(update)
		if !hasSentChatTask {
			bot.Send(msg)
		}
	}

}

func handleCommandMsg(update tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "start":
	case "chatgpt":
		msg.Text = "Hi, I'm ChatGPT bot. I can chat with you. Just send me a sentence and I will reply you. \n\n 请在这条消息下回复你的问题，我会回复你的。"
	case "ping":
		msg.Text = "pong"

	default:
		msg.Text = "I don't know that command"
	}
	return msg
}

func handleUserMessage(update tgbotapi.Update) (msg *tgbotapi.MessageConfig, hasSentChatTask bool) {
	_, thisUserHasMessage := session.Load(update.Message.From.ID)

	m := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg = &m
	hasSentChatTask = false

	if shouldIgnoreMsg(update) {
		return
	}

	if shouldLimitUser(update) {
		sendLimitMessage(update.Message.Chat.ID)
		return
	}
	if strings.Trim(update.Message.Text, " ") != "" {
		if !thisUserHasMessage {
			sendTaskToChannel(update.Message.Text, update.Message.Chat.ID, update.Message.From.ID, update.Message.MessageID)
			hasSentChatTask = true
		} else {
			log.Printf("[RateLimit] user %d is chatting with me, ignore message %s", update.Message.From.ID, update.Message.Text)
			sendRateLimitMessage(update.Message.Chat.ID)
		}
	} else {
		msg.Text = "Please provide a sentence."
	}

	return msg, hasSentChatTask
}

func sendLimitMessage(i int64) {
	text := fmt.Sprintf("You should join channel %s and group %s, then you can pm me", tgChannelName, tgGroupName) +
		"\n\n" + fmt.Sprintf("你需要加入频道 %s 和群组 %s，然后才能和我私聊", tgChannelName, tgGroupName)
	msg := tgbotapi.NewMessage(i, text)
	bot.Send(msg)
}

func findMemberFromChat(chatName string, userID int64) bool {
	findUserConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			SuperGroupUsername: chatName,
			UserID:             userID,
		},
	}
	member, err := bot.GetChatMember(findUserConfig)
	if err != nil || member.Status == "left" || member.Status == "kicked" {
		log.Printf("[ShouldLimitUser] memeber should be limit. id: %d", userID)
		return false
	}
	return true
}

func shouldLimitUser(update tgbotapi.Update) bool {
	// come from private
	if update.Message.Chat.IsPrivate() {
		userID := update.Message.From.ID
		canFindInChannel := findMemberFromChat(tgChannelName, userID)
		canFindInGroup := findMemberFromChat(tgGroupName, userID)
		if !(canFindInChannel && canFindInGroup) {
			return true
		}
	}
	return false
}

func shouldIgnoreMsg(update tgbotapi.Update) bool {
	// ignore message target not to me
	if update.Message == nil {
		return true
	}

	if update.Message.NewChatMembers != nil ||
		update.Message.LeftChatMember != nil {
		return true
	}

	return update.Message.ReplyToMessage != nil &&
		!update.Message.ReplyToMessage.From.IsBot
}

func sendRateLimitMessage(chat int64) {
	bot.Send(tgbotapi.NewMessage(chat, "you are chatting with me, please wait for a while."))
}

func loopAndFinishChatTask() {
	for {
		select {
		case task := <-TaskChannel:
			log.Println("[LoopAndFinishChatTask] got a task to finish")
			task.Finish()
		case <-time.After(30 * time.Second):
		}

	}
}
