package telegram

import (
	"chatgpt-bot/constant"
	"chatgpt-bot/constant/cmd"
	"chatgpt-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

type BotCmd = string

type CommandHandler interface {
	Cmd() BotCmd
	Run(b *Bot, update tgbotapi.Update) error
}

type StartCommand struct {
}

func (c *StartCommand) Cmd() BotCmd {
	return cmd.START
}

func (c *StartCommand) Run(b *Bot, update tgbotapi.Update) error {
	log.Println(fmt.Printf("get args: [%s]", update.Message.CommandArguments()))
	b.safeSendMsg(update.Message.Chat.ID, constant.BotStartTip)
	return nil
}

type ChatCommand struct {
}

func (c *ChatCommand) Cmd() BotCmd {
	return cmd.CHATGPT
}

func (c *ChatCommand) Run(b *Bot, update tgbotapi.Update) error {
	log.Println(fmt.Printf("get args: [%s]", update.Message.CommandArguments()))
	// todo: 处理邀请码
	b.safeSendMsg(update.Message.Chat.ID, constant.BotStartTip)
	return nil
}

type PingCommand struct {
}

func (c *PingCommand) Cmd() BotCmd {
	return cmd.PING
}

func (c *PingCommand) Run(b *Bot, update tgbotapi.Update) error {
	b.safeSendMsg(update.Message.Chat.ID, constant.BotPingTip)
	return nil
}

type LimiterCommand struct {
}

func (c *LimiterCommand) Cmd() BotCmd {
	return cmd.LIMITER
}

func (c *LimiterCommand) Run(b *Bot, update tgbotapi.Update) error {
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

type PprofCommand struct {
}

func (c *PprofCommand) Cmd() BotCmd {
	return cmd.PPROF
}

func (c *PprofCommand) Run(b *Bot, update tgbotapi.Update) error {
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

func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func NewLimiterCommand() *LimiterCommand {
	return &LimiterCommand{}
}

func NewPprofCommand() *PprofCommand {
	return &PprofCommand{}
}
