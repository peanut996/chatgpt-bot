package engine

import (
	"chatgpt-bot/cfg"
	"net/http"
)

type BingEngine struct {
	client *http.Client
}

func NewBingEngine() *BingEngine {
	return &BingEngine{
		client: &http.Client{},
	}
}

func (e *BingEngine) Init(cfg *cfg.Config) error {
	panic("implement me")
}

func (e *BingEngine) Chat(message string) (string, error) {
	panic("implement me")
}

func (e *BingEngine) Alive() bool {
	panic("implement me")
}
