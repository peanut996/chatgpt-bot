package logic

import (
	"context"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

var (
	client gpt3.Client
)

func init() {
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	engine := os.Getenv("OPEN_AI_ENGINE")
	c := gpt3.NewClient(apiKey, gpt3.WithDefaultEngine(engine))
	client = c
}

func ChatWithAI(sentence string) string {

	ctx := context.Background()
	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt:    []string{sentence},
		MaxTokens: gpt3.IntPtr(30),
		Stop:      []string{"."},
		Echo:      true,
	})
	if err != nil {
		return "unknown error"
	}
	return resp.Choices[0].Text
}
