package handler

import (
	"chatgpt-bot/bot/telegram"
	"chatgpt-bot/constant/cmd"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type VIPCommandHandler struct {
	userRepository *repository.UserRepository
}

func (v *VIPCommandHandler) Cmd() BotCmd {
	return cmd.VIP
}

func (v *VIPCommandHandler) Run(t telegram.TelegramBot, message tgbotapi.Message) error {
	if !t.IsBotAdmin(message.From.ID) {
		t.SafeSendMsg(message.Chat.ID, tip.NotAdminTip)
		return nil
	}
	args := message.CommandArguments()
	if args == "" {
		t.SafeSend(tgbotapi.NewMessage(message.Chat.ID, botError.MissingRequiredConfig+" : id"))
		return nil
	}
	user, err := v.userRepository.GetByUserID(args)
	if err != nil {
		return err
	}
	if user == nil {
		t.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, "user not found")
		return nil
	}
	err = v.userRepository.UpdateUserToVIP(args)
	if err != nil {
		return err
	}

	chatID, _ := utils.StringToInt64(user.UserID)
	t.SafeSendMsg(chatID, tip.BecomeDonorTip)
	t.SafeReplyMsgWithoutPreview(message.Chat.ID, message.MessageID, "success")
	return nil
}

func NewVIPCommandHandler(userRepository *repository.UserRepository) *VIPCommandHandler {
	return &VIPCommandHandler{userRepository: userRepository}
}
