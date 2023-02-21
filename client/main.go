package main

import (
	"chatgpt-bot/logic"
	"chatgpt-bot/wechat"
	"flag"
	"log"

	"github.com/joho/godotenv"
)

var botType string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.StringVar(&botType, "t", "tg", "which bot type to run. (tg/wechat) default: tg")
	flag.Parse()
}

func runTGBot() {
	log.Println("[Main] Start ChatGPT3 bot...")
	log.Println("[Main] Start Fetching Updates...")
	logic.FetchUpdates()

	log.Println("[Main] process exited")

}

func runWeChatBot() {
	wechat.InitBot()
	wechat.RegisterMessageHandler()
	wechat.GetWechatBot().Block()
}

func main() {
	if botType == "wechat" {
		log.Println("[Main] Start WeChat Bot...")
		runWeChatBot()
	} else {
		log.Println("[Main] Start Telegram Bot...")
		runTGBot()
	}
}
