package engine

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/model"
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

func (e *BingEngine) Chat(ctx model.ChatContext) (string, error) {
	panic("implement me")
}

func (e *BingEngine) Alive() bool {
	panic("implement me")
}
