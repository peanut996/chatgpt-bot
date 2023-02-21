package engine

import "chatgpt-bot/cfg"

var (
	ENGINE_CHATGPT = "chatgpt"
	ENGINE_BING    = "bing"
)

type Engine interface {
	Init(*cfg.Config) error
	Chat(string) (string, error)
	Alive() bool
}

func GetEngine(engineType string) Engine {
	switch engineType {
	case ENGINE_BING:
		return NewBingEngine()
	case ENGINE_CHATGPT:
		return NewChatGPTEngine()
	default:
		return NewChatGPTEngine()
	}
}
