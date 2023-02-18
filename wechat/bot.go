package wechat

import (
	"chatgpt-bot/logic"
	"fmt"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

var (
	bot *openwechat.Bot
)

func InitBot() {
	// bot = openwechat.DefaultBot()
	bot = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式
}

func RegisterMessageHandler() {
	// 注册消息处理函数
	bot.MessageHandler = HandleWechatMessage
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)
}

func GetWechatBot() *openwechat.Bot {
	return bot
}

func HandleWechatMessage(msg *openwechat.Message) {

	var originString = "@GPT机器人\u2005"

	if msg.IsText() && msg.IsAt() && strings.Contains(msg.Content, "GPT") {
		var text = strings.Replace(msg.Content, originString, "", -1)
		var res = logic.GetChatGPTResponseWithRetry(text)
		msg.ReplyText(res)
	}
}
