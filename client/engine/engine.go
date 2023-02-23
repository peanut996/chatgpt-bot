package engine

import "chatgpt-bot/cfg"

var (
	CHATGPT = "chatgpt"
	BING    = "bing"
)

type Engine interface {
	Init(*cfg.Config) error
	Chat(string) (string, error)
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
