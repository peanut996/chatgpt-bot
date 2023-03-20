package telegram

import (
	"chatgpt-bot/constant/cmd"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/constant/tip"
	"chatgpt-bot/model/persist"
	"chatgpt-bot/repository"
	"chatgpt-bot/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

type BotCmd = string

type CommandHandler interface {
	Cmd() BotCmd
	Run(b *Bot, message tgbotapi.Message) error
}

type PushDonateCommandHandler struct {
	userRepository *repository.UserRepository
}

func (p *PushDonateCommandHandler) Cmd() BotCmd {
	return cmd.PUSH
}

func (p *PushDonateCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	if !b.isBotAdmin(message.From.ID) {
		return fmt.Errorf(tip.NotAdminTip)
	}

	userIDs, err := p.userRepository.GetAllUserID()
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		go func(userID string) {
			if utils.IsEmpty(userID) {
				return
			}
			uid, _ := utils.StringToInt64(userID)
			msg := tgbotapi.NewMessage(uid, tip.DonateTip)
			b.safeSend(msg)
		}(userID)
	}

	return nil
}

type DonateCommandHandler struct{}

func (d *DonateCommandHandler) Cmd() BotCmd {
	return cmd.DONATE
}

func (d *DonateCommandHandler) Run(bot *Bot, message tgbotapi.Message) error {

	msg := tgbotapi.NewMessage(message.Chat.ID, tip.DonateTip)

	bot.safeSend(msg)

	return nil
}

type QueryCommandHandler struct {
	userRepository             *repository.UserRepository
	userInviteRecordRepository *repository.UserInviteRecordRepository
}

func (q *QueryCommandHandler) Cmd() BotCmd {
	return cmd.QUERY
}

func (q *QueryCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	userID := utils.Int64ToString(message.From.ID)
	user, err := q.userRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("[QueryCommandHandler] get user by user id failed, err: 【%s】\n", err)
		return err
	}
	if user == nil {
		userInfo, err := b.getUserInfo(message.From.ID)
		if err != nil {
			return err
		}
		err = q.userRepository.InitUser(userID, userInfo.String())
		if err != nil {
			log.Printf("[QueryCommandHandler] init user failed, err: 【%s】\n", err)
			return err
		}
		user, err = q.userRepository.GetByUserID(userID)
		if err != nil {
			log.Printf("[QueryCommandHandler] get user by user id failed, err: 【%s】\n", err)
			return err
		}
	}
	inviteCount, err := q.userInviteRecordRepository.CountByUserID(userID)
	if err != nil {
		log.Printf("[QueryCommandHandler] get user invite count by user id failed, err: 【%s】\n", err)
		return err
	}

	text := fmt.Sprintf(tip.QueryUserInfoTemplate,
		userID, user.RemainCount, inviteCount, b.getBotInviteLink(user.InviteCode))
	b.safeReplyMsg(message.Chat.ID, message.MessageID, text)
	return nil
}

func NewQueryCommandHandler(userRepository *repository.UserRepository, userInviteRecordRepository *repository.UserInviteRecordRepository) *QueryCommandHandler {
	return &QueryCommandHandler{
		userRepository,
		userInviteRecordRepository,
	}
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

func (c *StartCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	log.Println(fmt.Printf("get args: [%s]", message.CommandArguments()))
	args := message.CommandArguments()
	if matchInviteCode(args) {
		err := c.handleInvitation(args, utils.Int64ToString(message.From.ID), b)
		if err != nil {
			log.Printf("[StartCommandHandler] handle invitation failed, err: 【%s】", err)
		}
	}
	b.safeSendMsg(message.Chat.ID, tip.BotStartTip)
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
	if user.UserID == inviteUserID {
		log.Printf("[handleInvitation] user can not invite himself")
		return fmt.Errorf("[handleInvitation] user can not invite himself, user id: [%s]", inviteUserID)
	}
	record, err := c.userInviteRecordRepository.GetByInviteUserID(inviteUserID)
	if err != nil {
		log.Printf("[handleInvitation] find user by invite user id failed, err: 【%s】", err)
		return err
	}
	if record != nil {
		log.Printf("[handleInvitation]  user has been invited by other user: " + record.UserID)
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
	b.safeSendMsg(originUserID, tip.InviteSuccessTip)
	return nil
}

type ChatCommandHandler struct {
}

func (c *ChatCommandHandler) Cmd() BotCmd {
	return cmd.CHATGPT
}

func (c *ChatCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	log.Println(fmt.Printf("get args: [%s]", message.CommandArguments()))
	b.safeSendMsg(message.Chat.ID, tip.BotStartTip)
	return nil
}

type PingCommandHandler struct {
}

func (c *PingCommandHandler) Cmd() BotCmd {
	return cmd.PING
}

func (c *PingCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	b.safeSendMsg(message.Chat.ID, tip.BotPingTip)
	return nil
}

type LimiterCommandHandler struct {
}

func (c *LimiterCommandHandler) Cmd() BotCmd {
	return cmd.LIMITER
}

func (c *LimiterCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	if !b.isBotAdmin(message.From.ID) {
		msg.Text = tip.NotAdminTip
	} else {
		b.enableLimiter = utils.ParseBoolString(message.CommandArguments())
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

func (c *PprofCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	if !b.isBotAdmin(message.From.ID) {
		msg.Text = tip.NotAdminTip
		b.safeSend(msg)
		return nil
	}

	if filePath, success := dumpProfile(); success {
		defer func() {
			_ = os.Remove(filePath)
		}()
		err := sendFile(b, message.Chat.ID, filePath)
		if err == nil {
			return nil
		}
	}

	msg.Text = botError.InternalError
	b.safeSend(msg)
	return nil
}

func dumpProfile() (string, bool) {
	fileName := fmt.Sprintf("%d.pprof", time.Now().Unix())
	filePath := os.TempDir() + string(os.PathSeparator) + fileName
	tmpFile, err := os.Create(filePath)
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
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

func (i *InviteCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	userID := utils.Int64ToString(message.From.ID)
	user, err := i.userRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("[InviteCommandHandler] find user by user id failed, err: 【%s】", err)
		return err
	}
	if user != nil {
		link := b.getBotInviteLink(user.InviteCode)
		b.safeSendMsg(message.Chat.ID, fmt.Sprintf(tip.InviteTipTemplate, link, link))
		return nil
	} else {
		userName := ""
		tgUser, err := b.getUserInfo(message.From.ID)
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
		b.safeSendMsg(message.Chat.ID, fmt.Sprintf(tip.InviteTipTemplate, link, link))
	}
	return nil
}

type CountCommandHandler struct {
	userRepository *repository.UserRepository
}

func (c *CountCommandHandler) Cmd() BotCmd {
	return cmd.COUNT
}

func (c *CountCommandHandler) Run(b *Bot, message tgbotapi.Message) error {
	if !b.isBotAdmin(message.From.ID) {
		b.safeSendMsg(message.Chat.ID, tip.NotAdminTip)
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
		b.safeSendMsg(message.Chat.ID, fmt.Sprintf("failed to set count. params: %s, err: %s", args, err.Error()))
		return nil
	}
	b.safeSendMsg(message.Chat.ID, "success")
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

func NewCountCommandHandler(userRepository *repository.UserRepository) *CountCommandHandler {
	return &CountCommandHandler{
		userRepository: userRepository,
	}
}

func NewChatCommandHandler() *ChatCommandHandler {
	return &ChatCommandHandler{}
}

func NewDonateCommandHandler() *DonateCommandHandler {
	return &DonateCommandHandler{}
}
