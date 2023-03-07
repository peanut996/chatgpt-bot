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
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

// Chat is the method to chat with ChatGPT engine
func (e *ChatGPTEngine) chat(sentence string, userID string) (string, error) {
	log.Println("[ChatGPT] send request to chatgpt, text: ", sentence)

	if !e.Alive() {
		return constant.ChatGPTEngineNotOnline, nil
	}

	encodeSentence := url.QueryEscape(sentence)
	e.client.Timeout = time.Duration(constant.ChatGPTTimeoutSeconds) * time.Second
	queryString := fmt.Sprintf("/chat?user_id=%s&sentence=%s", userID, encodeSentence)
	resp, err := e.client.Get(e.baseUrl + queryString)
	defer resp.Body.Close()
	if err != nil {
		log.Println("[ChatGPT] chatgpt engine error: ", err)
		return "", errors.New(constant.NetworkError)
	}
	if resp.StatusCode != 200 {
		log.Println("[ChatGPT] chatgpt engine fail, status code: ", resp.StatusCode)
		return "", errors.New(constant.ChatGPTError)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ChatGPT] chatgpt engine error: ", err)
		return "", errors.New(constant.InternalError)
	}
	data := make(map[string]interface{}, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("[ChatGPT] unmarshal chatgpt response error: %s, resp: %s\n",
			err, string(body))
		return "", errors.New(constant.InternalError)
	}
	log.Println("[ChatGPT] response from chatgpt: ", utils.ToJsonString(data))
	if msg, ok := data["message"].(string); ok && msg != "" {
		return msg, nil
	}
	if detail, ok := data["detail"].(string); ok && detail != "" {
		return "", errors.New(fmt.Sprintf(constant.ChatGPTErrorTemplate, detail))
	}
	return "", errors.New(constant.ChatGPTError)
}

func (e *ChatGPTEngine) Chat(sentence string, userID string) (string, error) {
	resp, err := e.chat(sentence, userID)

	isNetworkError := strings.Contains(resp, "SSLError") ||
		strings.Contains(resp, "RemoteDisconnected") ||
		strings.Contains(resp, "ConnectionResetError")
	if isNetworkError {
		return "", errors.New(constant.NetworkError)
	}

	if err == nil && resp != "" {
		return resp, nil
	}

	if err != nil {
		return "", err
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
