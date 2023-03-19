package engine

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/model"
)

var (
	CHATGPT = "chatgpt"
	BING    = "bing"
)

type Engine interface {
	Init(*cfg.Config) error
	Chat(ctx model.ChatContext) (string, error)
	Alive() bool
}

func GetEngine(engineType string) Engine {
	switch engineType {
	case BING:
		return NewBingEngine()
	case CHATGPT:
		return NewChatGPTEngine()
	default:
		return NewChatGPTEngine()
	}
}
