package constant

var (
	BotPingTip    = "pong"
	UnknownCmdTip = "Unknown command, please send /start to start a chat \n\n" +
		"🔥未知命令，请发送 /start 来开始聊天"
	BotStartTip = "Hi, I'm ChatGPT bot. I can chat with you. Just send me a sentence and I will reply you. \nYou can also add me to your own group\n\n" +
		"😊请在这条消息下回复你的问题，我会回复你的 \n\n🔥你也可以私聊我或者把我加到你的群组聊天"
	OnlyOneChatAtATime = "you are chatting with me, please wait for a while. \n\n" +
		"😅你已经发送了一条信息，请耐心等待"
	LimitUserMessageTemplate = "You should join channel %s and group %s, then you can talk to me \n\n" +
		"😢你需要加入频道 %s 和群组 %s，然后才能和我交谈"

	RateLimitMessageTemplate = "You can only send %d messages in %d min, please try later. \nRate limiter will disappeared when you join both channel %s and group %s\n\n" +
		"😅 你只能在 %d 分钟内发送 %d 条消息，请稍候再试\n" +
		"当你同时加入频道 %s 和群组 %s 后，将不再限速"

	ChatGPTError = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧"
	ChatGPTErrorTemplate = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧 \n\n %s"
	ChatGPTEngineNotOnline = "Chatgpt engine is not ready, please wait a moment. \n\n" +
		"😇ChatGPT 引擎还没有准备好，请稍等一下"
	SendBackMsgFailed = "Send back message failed, please try again later \n\n" +
		"😇返回消息失败，请稍后再试"

	NetworkError = "Network error, please try again later \n\n" +
		"😐网络错误，请稍后再试"
)
