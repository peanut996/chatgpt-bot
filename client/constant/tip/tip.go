package tip

import "fmt"

var (
	BotPingTip    = "pong"
	UnknownCmdTip = "Unknown command, please send /start to start a chat \n\n" +
		"🔥未知命令，请发送 /start 来开始聊天"
	BotStartTip = "Just send me a sentence and I will reply you. \n" +
		"You can also add me to your own group\n" +
		"Bot default use gpt-3.5 model, if you want to use gpt-4 model, please use `/gpt4` command, like 【/gpt4 how is weather today?】 \n\n" +
		"😊请在这条消息下回复你的问题，我会回复你的 \n" +
		"🔥你也可以私聊我或者把我加到你的群组聊天 \n" +
		"🤖默认使用gpt-3.5模型，gpt-4模型对话请使用「 /gpt4 」+ 空格 + 你的问题，如【/gpt4 今天天气怎么样?】"

	NotAdminTip = "You are not admin, can not use this command \n\n" +
		"😢你不是管理员，无法执行此操作"

	InviteSuccessTip = "Invite success, you can use /invite to get your invite link \n\n" +
		"😊邀请成功，你可以使用 /invite 来获取你的邀请链接"

	InviteTipTemplate = "You can invite 1 new user after that you can use gpt4 chat unlimited. Invite link: %s\n\n" +
		"😊你可以邀请1个新用户后可无限使用gpt4对话. 你的邀请链接: %s"

	QueryUserInfoTemplate = "💁账号(Account): %s\n\n" +
		"🏆剩余次数(RemainCount): %d\n" +
		"🎭邀请人数(InviteUsers): %d\n" +
		"🔗邀请链接(InviteLink): %s\n\n" +
		"🔮小提示：邀请1个新用户后可无限使用gpt4对话(Invite one new user to get gpt4 unlimited)"

	GPTLackTextTipTemplate = "`/%s` + blank + your question.\n\n" +
		"😊「 /%s 」+ 空格 + 你的问题"

	AlipayQRCodeUrl = "https://raw.githubusercontent.com/peanut996/chatgpt-bot/master/assets/alipay.JPG"

	WechatQRCodeUrl = "https://raw.githubusercontent.com/peanut996/chatgpt-bot/master/assets/wechat.JPG"

	DonateTip = fmt.Sprintf("🙏 感谢您使用我们的机器人！如果您觉得我们的机器人对您有所帮助，欢迎为我们捐赠，以支持我们的运营和发展。\n\n"+
		"💰 您可以通过以下方式向我们捐赠：\n\n- [微信](%s)\n\n- [支付宝](%s) \n\n"+
		"💡 如果您有任何其他的捐赠方式或者建议，欢迎联系我们！\n\n"+
		"👏 再次感谢您的支持，您的捐赠将帮助我们更好地为您提供服务！\n", WechatQRCodeUrl, AlipayQRCodeUrl)

	StatusTipTemplate = "💁 总用户数：%d\n\n" +
		"🏆 总邀请记录次数：%d\n"

	QueueTipTemplate = "Current queue count: %d, please wait wait patiently\n\n" +
		"🔥当前排队人数：%d, 请耐心等候\n\n"
)
