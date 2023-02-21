package engine

import (
	"chatgpt-bot/cfg"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ChatGPTEngine struct {
	client *http.Client

	baseUrl string

	alive bool
}

func NewChatGPTEngine() *ChatGPTEngine {
	return &ChatGPTEngine{}
}

func (e *ChatGPTEngine) Init(cfg *cfg.Config) error {
	e.client = &http.Client{}
	e.baseUrl = fmt.Sprintf("http://%s:%d", cfg.EngineConfig.Host, cfg.EngineConfig.Port)

	go e.checkChatGPTEngine()
	return nil
}

func (e *ChatGPTEngine) Alive() bool {
	resp, err := http.Get(e.baseUrl + "/ping")
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

// Chat is the method to chat with ChatGPT engine
func (e *ChatGPTEngine) chat(sentence string) (string, error) {
	log.Println("[ChatGPT] send request to chatgpt, text: ", sentence)

	if !e.Alive() {
		return "chatgpt engine is not ready, please wait a moment.", nil
	}

	encodeSentence := url.QueryEscape(sentence)
	e.client.Timeout = 300 * time.Second
	resp, err := e.client.Get(e.baseUrl + "/chat?sentence=" + encodeSentence)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("chatgpt engine return error")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println("[ChatGPT] response from chatgpt: ", string(body))
	data := make(map[string]string, 0)
	json.Unmarshal(body, &data)
	if data["message"] == "" {
		if data["detail"] != "" {
			return data["detail"], nil
		}
		return "", errors.New("chatgpt engine return empty, may be too many requests in one hour, try again later")
	}
	return data["message"], nil
}

func (e *ChatGPTEngine) Chat(sentence string) (string, error) {
	resp, err := e.chat(sentence)
	if err == nil {
		return resp, nil
	}
	return "", err
}

func (e *ChatGPTEngine) checkChatGPTEngine() {
	for {
		status := e.Alive()
		if !status {
			log.Println("[HealthCheck] chatgpt engine is not ready")
			e.alive = false
		} else {
			e.alive = true
		}
		time.Sleep(10 * time.Second)
	}
}
