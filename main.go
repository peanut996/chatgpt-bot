package main

import (
	"chatgpt-bot/logic"
	"fmt"
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
	log.Println("Start ChatGPT3 bot...")

	log.Println("Start fetching offset...")
	// logic.RefreshLastestOffset()
	log.Println("Fetching offset Completed..")

	log.Println("Start Fetching Updates...")
	logic.FetchUpdates()

	fmt.Println("process exited")
}
