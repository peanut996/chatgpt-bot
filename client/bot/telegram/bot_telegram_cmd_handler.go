package telegram

import (
	"chatgpt-bot/constant"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/model/persist"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

type BotCmd = string

type CommandHandler interface {
	Cmd() BotCmd
	Run(b *Bot, update tgbotapi.Update) error
}

type StartCommandHandler struct {
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
}

func (c *StartCommandHandler) Cmd() BotCmd {
	return cmd.START
}

func matchInviteCode(code string) bool {
	return utils.IsNotEmpty(code) && len(code) == 10 && utils.IsMatchString(`^[a-zA-Z]{10}$`, code)
}

func (c *StartCommandHandler) Run(b *Bot, update tgbotapi.Update) error {
	log.Println(fmt.Printf("get args: [%s]", update.Message.CommandArguments()))
	args := update.Message.CommandArguments()
	if matchInviteCode(args) {
		err := c.handleInvitation(args, utils.ConvertInt64ToString(update.Message.From.ID), b)
		if err != nil {
			log.Printf("[StartCommandHandler] handle invitation failed, err: 【%s】", err)
		}
	}
	b.safeSendMsg(update.Message.Chat.ID, constant.BotStartTip)
	return nil
}

func (c *StartCommandHandler) handleInvitation(inviteCode string, inviteUserID string, b *Bot) error {
	user, err := c.userRepository.GetUserByInviteCode(inviteCode)
	if err != nil {
		log.Printf("[handleInvitation] find user by invite code failed, err: 【%s】", err)
		return err
	}
	if user == nil {
		log.Printf("[handleInvitation] find user by invite code failed, user is nil")
		return errors.New("no such user by invite code: " + inviteCode)
	}
	inviteUser, err := c.userRepository.GetByUserID(inviteUserID)
	if err != nil {
		log.Printf("[handleInvitation] find user by invite user id failed, err: 【%s】", err)
		return err
	}
	if inviteUser != nil {
		log.Printf("[handleInvitation]  user has been invited by other user" + inviteUserID)
		return nil
	}
	inviteRecord := persist.NewUserInviteRecord(user.UserID, inviteUserID)
	err = c.userInviteRecordRepository.Insert(inviteRecord)
	if err != nil {
		return err
	}
	err = c.userRepository.AddCountWhenInviteOther(user.UserID)
	if err != nil {
		return err
	}
	originUserID, _ := utils.StringToInt64(user.UserID)
	b.safeSendMsg(originUserID, constant.InviteSuccessTip)
	return nil
}

type ChatCommandHandler struct {
}

func (c *ChatCommandHandler) Cmd() BotCmd {
	return cmd.CHATGPT
}

func (c *ChatCommandHandler) Run(b *Bot, update tgbotapi.Update) error {
	log.Println(fmt.Printf("get args: [%s]", update.Message.CommandArguments()))
	b.safeSendMsg(update.Message.Chat.ID, constant.BotStartTip)
	return nil
}

type PingCommandHandler struct {
}

func (c *PingCommandHandler) Cmd() BotCmd {
	return cmd.PING
}

func (c *PingCommandHandler) Run(b *Bot, update tgbotapi.Update) error {
	b.safeSendMsg(update.Message.Chat.ID, constant.BotPingTip)
	return nil
}

type LimiterCommandHandler struct {
}

func (c *LimiterCommandHandler) Cmd() BotCmd {
	return cmd.LIMITER
}

func (c *LimiterCommandHandler) Run(b *Bot, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if !b.isBotAdmin(update.Message.From.ID) {
		msg.Text = constant.NotAdminTip
	} else {
		b.enableLimiter = utils.ParseBoolString(update.Message.CommandArguments())
		msg.Text = fmt.Sprintf("limiter status is %v now", b.enableLimiter)
	}
	b.safeSend(msg)
	return nil
}

type PprofCommandHandler struct {
}

func (c *PprofCommandHandler) Cmd() BotCmd {
	return cmd.PPROF
}

func (c *PprofCommandHandler) Run(b *Bot, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if !b.isBotAdmin(update.Message.From.ID) {
		msg.Text = constant.NotAdminTip
		b.safeSend(msg)
		return nil
	}

	if filePath, success := dumpProfile(); success {
		err := sendFile(b, update.Message.Chat.ID, filePath)
		if err == nil {
			return nil
		}
	}

	msg.Text = constant.InternalError
	b.safeSend(msg)
	return nil
}

func dumpProfile() (string, bool) {
	fileName := fmt.Sprintf("%d.pprof", time.Now().Unix())
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	tmpFile, err := os.Create(filePath)
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
		_ = os.Remove(filePath)
	}(tmpFile)

	if err != nil {
		log.Printf("[DumpProfile] create temp file failed, err: 【%s】", err)
		return err.Error(), false
	}

	err = pprof.WriteHeapProfile(tmpFile)
	if err != nil {
		log.Printf("[DumpProfile] create temp file failed, err: 【%s】", err)
		return err.Error(), false
	}

	return tmpFile.Name(), true
}

func sendFile(b *Bot, chatID int64, filePath string) error {
	fileMsg := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filePath))
	_, err := b.tgBot.Send(fileMsg)
	if err != nil {
		log.Printf("[SendFile] send file failed, err: 【%s】", err)
		return err
	}
	return nil
}

type InviteCommandHandler struct {
	userRepository *repository.UserRepository
}

func (i *InviteCommandHandler) Cmd() BotCmd {
	return cmd.INVITE
}

func (i *InviteCommandHandler) Run(b *Bot, update tgbotapi.Update) error {
	userID := utils.ConvertInt64ToString(update.Message.From.ID)
	user, err := i.userRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("[InviteCommandHandler] find user by user id failed, err: 【%s】", err)
		return err
	}
	if user != nil {
		link := b.getBotInviteLink(user.InviteCode)
		b.safeSendMsg(update.Message.Chat.ID, fmt.Sprintf(constant.InviteTipTemplate, link, link))
		return nil
	} else {
		userName := ""
		tgUser, err := b.getUserInfo(update.Message.From.ID)
		if err == nil {
			userName = tgUser.String()
		}
		err = i.userRepository.InitUser(userID, userName)
		if err != nil {
			log.Printf("[InviteCommandHandler] init user failed, err: 【%s】", err)
			return err
		}
		user, _ := i.userRepository.GetByUserID(userID)
		link := b.getBotInviteLink(user.InviteCode)
		b.safeSendMsg(update.Message.Chat.ID, fmt.Sprintf(constant.InviteTipTemplate, link, link))
	}
	return nil
}

func NewStartCommandHandler(userRepository *repository.UserRepository, userInviteRecordRepository *repository.UserInviteRecordRepository) *StartCommandHandler {
	return &StartCommandHandler{
		userRepository:             userRepository,
		userInviteRecordRepository: userInviteRecordRepository,
	}
}

func NewPingCommandHandler() *PingCommandHandler {
	return &PingCommandHandler{}
}

func NewLimiterCommandHandler() *LimiterCommandHandler {
	return &LimiterCommandHandler{}
}

func NewPprofCommandHandler() *PprofCommandHandler {
	return &PprofCommandHandler{}
}

func NewInviteCommandHandler(userRepository *repository.UserRepository) *InviteCommandHandler {
	return &InviteCommandHandler{
		userRepository: userRepository,
	}
}
