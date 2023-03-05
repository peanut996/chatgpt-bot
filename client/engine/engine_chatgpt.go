package engine

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/constant"
	"chatgpt-bot/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
		return constant.ChatGPTEngineNotOnline, nil
	}

	encodeSentence := url.QueryEscape(sentence)
	e.client.Timeout = time.Duration(constant.ChatGPTTimeoutSeconds) * time.Second
	resp, err := e.client.Get(e.baseUrl + "/chat?sentence=" + encodeSentence)
	if err != nil {
		log.Println(err)
		return "", errors.New(constant.NetworkError)
	}
	if resp.StatusCode != 200 {
		return "", errors.New(constant.ChatGPTError)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", errors.New(constant.InternalError)
	}
	data := make(map[string]string, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	log.Println("[ChatGPT] response from chatgpt: ", utils.ToJsonString(data))
	if data["message"] == "" {
		if data["detail"] != "" {
			return data["detail"], nil
		}
		return "", errors.New(constant.ChatGPTError)
	}
	return data["message"], nil
}

func (e *ChatGPTEngine) Chat(sentence string) (string, error) {
	resp, err := e.chat(sentence)

	isNetworkError := strings.Contains(resp, "SSLError") ||
		strings.Contains(resp, "RemoteDisconnected") ||
		strings.Contains(resp, "ConnectionResetError")
	if isNetworkError {
		return "", errors.New(constant.NetworkError)
	}

	if err == nil && "" != resp {
		return resp, nil
	}

	if err != nil {
		log.Println("[ChatGPT] chatgpt engine error: ", err)
		return fmt.Sprintf(constant.ChatGPTErrorTemplate, err.Error()), nil
	}
	return constant.ChatGPTError, err
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
