package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

type CountCommandHandler struct {
	userRepository *repository.UserRepository
}

func (c *CountCommandHandler) Cmd() BotCmd {
	return cmd.COUNT
}

func (c *CountCommandHandler) Run(b telegram.TelegramBot, message tgbotapi.Message) error {
	if !b.IsBotAdmin(message.From.ID) {
		b.SafeSendMsg(message.Chat.ID, tip.NotAdminTip)
		return nil
	}
	args := message.CommandArguments()
	if args == "" {
		return fmt.Errorf("invalid args")
	}
	params := strings.Split(args, ":")
	if len(params) != 2 {
		return fmt.Errorf("invalid args")
	}
	err := c.userRepository.UpdateCountByUserID(params[0], params[1])
	if err != nil {
		log.Printf("failed to set count. params: %s, err: %s", args, err.Error())
		b.SafeSendMsg(message.Chat.ID, fmt.Sprintf("failed to set count. params: %s, err: %s", args, err.Error()))
		return nil
	}
	b.SafeSendMsg(message.Chat.ID, "success")
	return nil
}

func NewCountCommandHandler(userRepository *repository.UserRepository) *CountCommandHandler {
	return &CountCommandHandler{
		userRepository: userRepository,
	}
}
