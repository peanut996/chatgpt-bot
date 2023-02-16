package logic

import (
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	chatId int64
	bot    *tgbotapi.BotAPI
)

func init() {
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
	bot.Send(tgbotapi.NewMessage(chatId, sentence))
	return sentence
}
