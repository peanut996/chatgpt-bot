package logic

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	chatId int64
	bot    *tgbotapi.BotAPI
	offset int = 0
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
	log.Printf("[%d] From [%s]: %s", update.UpdateID, update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "ping":
			msg.Text = "pong"
		default:
			msg.Text = "I don't know that command"
		}
	} else {
		if update.Message.Text != "" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, ChatWithAI(update.Message.Text))
		}
	}
	bot.Send(msg)
}

func RefreshLastestOffset() {
	o := 0

	u := tgbotapi.UpdateConfig{}
	u.Timeout = 60

	updates, err := bot.GetUpdates(u)
	if err != nil {
		log.Panic(err)
	}

	for _, update := range updates {
		if update.UpdateID >= o {
			o = update.UpdateID + 1
		}
	}
	fmt.Println("Starting from offset " + strconv.Itoa(o))
	offset = o
}
