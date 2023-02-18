package main

import (
	"chatgpt-bot/wechat"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// log.Println("[Main] Start ChatGPT3 bot...")
	// log.Println("[Main] Start Fetching Updates...")
	// logic.FetchUpdates()

	// log.Println("[Main] process exited")

	wechat.InitBot()
	wechat.RegisterMessageHandler()
	wechat.GetWechatBot().Block()
}
