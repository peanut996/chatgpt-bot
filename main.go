package main

import (
	"chatgpt-bot/api"
	"fmt"
	"log"

	gin "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	fmt.Println("hello world")
	r := gin.Default()
	r.GET("/chat", api.Chat)
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
