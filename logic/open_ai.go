package logic

import (
	"context"
	"log"
	"os"
	"time"

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
	var resp *gpt3.CompletionResponse
	var err error
	for i := 0; i < 10; i++ {
		log.Println("send request to open ai, text: ", sentence)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
		defer cancel()
		resp, err = client.Completion(ctx, gpt3.CompletionRequest{
			Prompt:    []string{sentence},
			MaxTokens: gpt3.IntPtr(4000),
			Echo:      false,
		})
		if err != nil {
			log.Println(err)
			continue
		} else {
			return resp.Choices[0].Text
		}
	}
	return "exceed max retry."
}
