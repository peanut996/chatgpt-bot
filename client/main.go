package main

import (
	"chatgpt-bot/app"
	"chatgpt-bot/cfg"
)

func main() {
	c, err := cfg.InitConfig()
	if err != nil {
		panic(err)
	}
	a := app.GetApp()

	err = a.Init(c)
	if err != nil {
		panic(err)
	}
	a.Run()
	a.Block()
}
