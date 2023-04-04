package tip

import (
	"chatgpt-bot/constant/config"
	"fmt"
)

var (
	BotPingTip    = "pong"
	UnknownCmdTip = "Unknown command, please send /start to start a chat \n\n" +
		"🔥未知命令，请发送 /start 来开始聊天"
	BotStartTip = "Please reply to this message with your question, and I will respond to you. \n" +
		"You can also private message me or add me to your group chat. \n" +
		"If you want to chat with me, use the command '/gpt' followed by a space and your question, for example, '/gpt How is the weather today?'. For conversations with GPT-4, please use the command '/gpt4'.\n\n" +
		"😊请在这条消息下回复你的问题，我会回复你的 \n" +
		"🔥你也可以私聊我或者把我加到你的群组聊天 \n" +
		"🤖如果你想和我聊些什么「 /gpt 」+ 空格 + 你的问题，如【/gpt 今天天气怎么样?】. GPT-4对话请使用「 /gpt4 」"

	NotAdminTip = "You are not admin, can not use this command \n\n" +
		"😢你不是管理员，无法执行此操作"

	InviteSuccessTip = "Invite success, you can use /invite to get your invite link \n\n" +
		"😊邀请成功，你可以使用 /invite 来获取你的邀请链接"

	InviteTipTemplate = "You can invite %d new user after that you can use /gpt4 chat unlimited. Invite link: %s. You can still use /gpt chat.\n\n" +
		"😊邀请%d个新用户后可使用gpt4对话. 你的邀请链接: %s\n\n" +
		"🔮小提示：「 /gpt 」命令没有邀请人数限制"

	InviteLinkTemplate = "Invite link: %s.\n\n" +
		"😊你的邀请链接: %s\n\n"

	QueryUserInfoTemplate = "💁账号(Account): %s\n\n" +
		"🥇捐赠用户(Donated): %t\n" +
		"🏆剩余次数(RemainCount): %d\n" +
		"🎭邀请人数(InviteUsers): %d\n" +
		"🔗邀请链接(InviteLink): %s\n\n" +
		"🔮小提示：成为捐赠用户或邀请%d个新用户后可使用gpt4对话(Invite %d new user to get gpt4 unlimited)"

	GPTLackTextTipTemplate = "`/%s` + blank + your question.\n\n" +
		"😊「 /%s 」+ 空格 + 你的问题"

	DonateTip = fmt.Sprintf("🙏 感谢您使用我们的机器人！如果您觉得我们的机器人对您有所帮助，欢迎为我们捐赠，以支持我们的运营和发展。\n\n"+
		"💰 您可以通过以下方式向我们捐赠：\n\n- [微信(Wechat)](%s)\n\n- [支付宝(Alipay)](%s) \n\n"+
		"💡 如果您有任何其他的捐赠方式或者建议，欢迎联系我们！\n\n"+
		"👏 再次感谢您的支持，您的捐赠将帮助我们更好地为您提供服务！\n\n"+
		"🔮 提示: 捐赠用户不再会接收此提示!", config.WechatQRCodeUrl, config.AlipayQRCodeUrl)

	StatusTipTemplate = "💁 总用户数：%d\n\n" +
		"🏆 总邀请记录次数：%d\n"

	QueueTipTemplate = "Current queue count: %d, please wait wait patiently\n\n" +
		"🔥当前排队人数：%d, 请耐心等候\n\n"

	AccessCodeTipTemplate = "Your access code: `%s`\n\n" +
		"😊你的访问码是: `%s`\n\n"
)
