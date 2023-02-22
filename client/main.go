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
	app := app.GetApp()

	err = app.Init(c)
	if err != nil {
		panic(err)
	}
	app.Run()
	app.Block()
}
