package logic

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	chatId int64
	bot    *tgbotapi.BotAPI
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
