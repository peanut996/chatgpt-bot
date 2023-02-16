package logic

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

var (
	client gpt3.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	engine := os.Getenv("OPEN_AI_ENGINE")
	c := gpt3.NewClient(apiKey, gpt3.WithDefaultEngine(engine))
	client = c
}

func ChatWithAI(sentence string) string {
	ctx := context.Background()
	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt:    []string{sentence},
		MaxTokens: gpt3.IntPtr(4000),
		Echo:      false,
	})
	if err != nil {
		return "unknown error"
	}
	fmt.Println(resp.Choices[0].Text)
	return resp.Choices[0].Text
}
