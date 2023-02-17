package logic

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

var (
	client        gpt3.Client
	isEngineReady bool = false
	engineUrl     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	engine := os.Getenv("OPEN_AI_ENGINE")

	if os.Getenv("ENGINE_URL") != "" {
		engineUrl = os.Getenv("ENGINE_URL")
	} else {
		engineUrl = "http://127.0.0.1:5000"
	}

	c := gpt3.NewClient(apiKey, gpt3.WithDefaultEngine(engine))
	client = c

	go checkChatGPTEngine()
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

func ChatWithChatGPT(sentence string) (string, error) {
	log.Println("[ChatGPT] send request to chatgpt, text: ", sentence)

	if !isEngineReady {
		return "chatgpt engine is not ready, please wait a moment.", nil
	}
	// encode sentence
	encodeSentence := url.QueryEscape(sentence)
	resp, err := http.Get(engineUrl + "/chat?sentence=" + encodeSentence)
	if err != nil {
		log.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println("[ChatGPT] response from chatgpt: ", string(body))
	data := make(map[string]string, 0)
	json.Unmarshal(body, &data)
	return data["message"], nil
}

func GetChatGPTResponseWithRetry(sentence string) string {
	for i := 0; i < 10; i++ {
		resp, err := ChatWithChatGPT(sentence)
		if err == nil {
			return resp
		}
	}
	return "Get gpt bot response fail, exceed max retry 10 times."
}

func healthCheck() bool {
	resp, err := http.Get(engineUrl + "/ping")
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

func checkChatGPTEngine() {
	for {
		status := healthCheck()
		if !status {
			log.Println("[HealthCheck] chatgpt engine is not ready")
			isEngineReady = false
		} else {
			isEngineReady = true
			log.Println("[HealthCheck] chatgpt engine is ready")
		}
		time.Sleep(10 * time.Second)
	}
}

func GetIsEngineReady() bool {
	return isEngineReady
}
