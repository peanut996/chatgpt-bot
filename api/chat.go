package api

import (
	"chatgpt-bot/logic"

	"github.com/gin-gonic/gin"
)

func Chat(c *gin.Context) {
	sentence := c.Query("sentence")
	if sentence == "" {
		c.String(400, "sentence is empty")
		return
	}
	logic.SendMessageToBot(sentence)
	c.String(200, sentence)
}
