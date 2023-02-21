package bot

import (
	"chatgpt-bot/cfg"
	"chatgpt-bot/engine"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
)

type WechatBot struct {
	engine engine.Engine
	bot    *openwechat.Bot

	botName string
}

var (
	NORMAL  = "normal"
	DESKTOP = "desktop"

	tag = "\u2005"
)

func (w *WechatBot) Init(cfg *cfg.Config) error {
	switch cfg.BotConfig.WechatLoginType {
	case DESKTOP:
		w.bot = openwechat.DefaultBot(openwechat.Desktop)
	case NORMAL:
		w.bot = openwechat.DefaultBot(openwechat.Normal)
	default:
		w.bot = openwechat.DefaultBot(openwechat.Normal)
	}
	w.Register()

	w.engine = engine.GetEngine(cfg.EngineConfig.EngineType)

	err := w.engine.Init(cfg)
	if err != nil {
		return err
	}

	w.botName = cfg.BotConfig.WechatBotName

	log.Printf("[Init] wechat bot init success, bot name: %s", w.botName)
	return nil
}

func NewWechatBot() *WechatBot {
	return &WechatBot{}
}

func (w *WechatBot) Register() {

	w.bot.MessageHandler = w.handleWechatMessage

	w.bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

}

func (w *WechatBot) Run() {
	if err := w.bot.Login(); err != nil {
		fmt.Println(err)
		return
	}
	self, err := w.bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	friends, err := self.Friends()
	fmt.Println(friends, err)

	groups, err := self.Groups()
	fmt.Println(groups, err)
	go w.loopAndCheckAlive()

	log.Println("[Run] wechat bot is running")
}

func (w *WechatBot) loopAndCheckAlive() {
	for {
		if !w.bot.Alive() {
			log.Println("wechat bot is not alive, try to login again")
		}
		time.Sleep(10 * time.Second)
	}
}
func (w *WechatBot) handleWechatMessage(msg *openwechat.Message) {
	atTag := fmt.Sprintf("@%s%s", w.botName, tag)

	if msg.IsText() && msg.IsSendByGroup() && msg.IsAt() && strings.Contains(msg.Content, w.botName) {
		sender, err := msg.SenderInGroup()
		if err != nil {
			return
		}
		var text = strings.Replace(msg.Content, atTag, "", -1)
		res, err := w.engine.Chat(text)

		replyText := ""
		if err != nil {
			log.Println(err)
			replyText = fmt.Sprintf("@%s%s", sender.NickName, tag) + "\n" + err.Error()
		} else {
			replyText = fmt.Sprintf("@%s%s", sender.NickName, tag) + "\n" + res
		}
		msg.ReplyText(replyText)
		msg.AppMsgType = openwechat.AppMsgTypeAttach
	}
}
