package logic

import (
	"context"
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
	resp, err := client.Completion(context.Background(), gpt3.CompletionRequest{
		Prompt:    []string{sentence},
		MaxTokens: gpt3.IntPtr(4000),
		Echo:      false,
	})
	if err != nil {
		return err.Error()
	}
	return resp.Choices[0].Text
}
