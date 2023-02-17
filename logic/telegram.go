package logic

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type ChatTask struct {
	Question string
	Answer   string
	Chat     int64
	From     int64
}

var (
	chatId      int64
	bot         *tgbotapi.BotAPI
	offset      int = 0
	session     *sync.Map
	TaskChannel chan *ChatTask
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	id, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	b, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}
	bot = b
	chatId = id
	session = &sync.Map{}
	TaskChannel = make(chan *ChatTask, 1)

	go loopAndFinishChatTask()
}

func NewChatTask(question string, chat, from int64) *ChatTask {
	return &ChatTask{
		Question: question,
		Chat:     chat,
		From:     from,
	}
}

func SendTaskToChannel(question string, chat, from int64) {
	log.Printf("[SendTaskToChannel] with question %s, chat id: %d, from: %d", question, chat, from)
	chatTask := NewChatTask(question, chat, from)
	TaskChannel <- chatTask
}

func (t *ChatTask) Send() {
	msg := tgbotapi.NewMessage(t.Chat, t.Question)
	msg.ParseMode = "markdown"
	msg.Text = t.Answer
	bot.Send(msg)
}

func (t *ChatTask) GetAnswerFromChatGPT() {
	a := GetChatGPTResponseWithRetry(t.Question)
	t.Answer = a
}

func (t *ChatTask) GetAnswerWithGPT() {
	t.Answer = ChatWithAI(t.Question)
}

func (t *ChatTask) Finish() {
	log.Printf("[Finish] start chat task with question %s, chat id: %d, from: %d", t.Question, t.Chat, t.From)
	defer session.Delete(t.From)

	t.GetAnswerFromChatGPT()
	t.Send()

	log.Printf("[Finish] end chat task with question %s, chat id: %d, from: %d", t.Question, t.Chat, t.From)

}

func SendMessageToBot(sentence string) string {
	response := ChatWithAI(sentence)
	_, err := bot.Send(tgbotapi.NewMessage(chatId, response))
	if err != nil {
		return err.Error()
	}
	return response
}

func FetchUpdates() {
	config := tgbotapi.NewUpdate(offset)
	config.Timeout = 60

	for update := range bot.GetUpdatesChan(config) {
		go handleUpdate(update)
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[BotUpdate] update id:[%d] from [%s] : %s", update.UpdateID, update.Message.From.String(), update.Message.Text)

	from := update.Message.From.ID

	_, thisUserHasMessage := session.Load(from)

	if thisUserHasMessage {
		log.Println("[MessageLimiter] this user has message, ignore this message. from: ", update.Message.From.ID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = "you are chatting with me, please wait for a while."
		bot.Send(msg)
		return
	} else {
		session.Store(from, &struct{}{})
		msg, hasSentChatTask := handleUserMessage(update)
		if !hasSentChatTask {
			bot.Send(msg)
			defer session.Delete(from)
		}
	}
}

func handleUserMessage(update tgbotapi.Update) (msg *tgbotapi.MessageConfig, hasSentChatTask bool) {
	m := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg = &m
	hasSentChatTask = false
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "ping":
			msg.Text = "pong"
		case "chat":
			if strings.Trim(update.Message.CommandArguments(), " ") != "" {
				SendTaskToChannel(update.Message.CommandArguments(), update.Message.Chat.ID, update.Message.From.ID)
				hasSentChatTask = true
			} else {
				msg.Text = "Please provide a sentence."
				return msg, false
			}
		default:
			msg.Text = "I don't know that command"
		}
	} else {
		if strings.Trim(update.Message.Text, " ") != "" {
			SendTaskToChannel(update.Message.Text, update.Message.Chat.ID, update.Message.From.ID)
			hasSentChatTask = true
		} else {
			msg.Text = "Please provide a sentence."
		}
	}
	return msg, hasSentChatTask
}

func loopAndFinishChatTask() {
	for {
		select {
		case task := <-TaskChannel:
			log.Println("[LoopAndFinishChatTask] got a task to finish")
			task.Finish()
		case <-time.After(30 * time.Second):
			log.Println("[LoopAndFinishChatTask] timeout after 30s")
		}

	}
}
