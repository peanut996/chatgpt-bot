package logic

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	isEngineReady bool = false
	engineUrl     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("ENGINE_URL") != "" {
		engineUrl = os.Getenv("ENGINE_URL")
	} else {
		engineUrl = "http://127.0.0.1:5000"
	}

	go checkChatGPTEngine()
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
	if resp.StatusCode != 200 {
		return "", errors.New("chatgpt engine return error")
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
		}
		time.Sleep(10 * time.Second)
	}
}

func GetIsEngineReady() bool {
	return isEngineReady
}
