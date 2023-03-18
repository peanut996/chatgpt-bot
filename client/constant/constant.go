package constant

var (
	ChatGPTTimeoutSeconds = 360

	DefaultCount = 10

	CountWhenInviteOtherUser = 30
)

var (
	BotPingTip    = "pong"
	UnknownCmdTip = "Unknown command, please send /start to start a chat \n\n" +
		"🔥未知命令，请发送 /start 来开始聊天"
	BotStartTip = "Hi, I'm ChatGPT bot. I can chat with you. Just send me a sentence and I will reply you. \nYou can also add me to your own group\n\n" +
		"😊请在这条消息下回复你的问题，我会回复你的 \n\n🔥你也可以私聊我或者把我加到你的群组聊天"

	NotAdminTip = "You are not admin, can not use this command \n\n" +
		"😢你不是管理员，无法执行此操作"

	InviteSuccessTip = "Invite success, you can use /invite to get your invite link \n\n" +
		"😊邀请成功，你可以使用 /invite 来获取你的邀请链接"

	InviteTipTemplate = "You can invite new users to get 50 chat sessions per new user. your invite link: %s\n\n" +
		"😊你可以邀请新用户获取聊天次数 30次/新用户. 你的邀请链接: %s"
)

var (
	OnlyOneChatAtATime = "you are chatting with me, please wait for a while. \n\n" +
		"😅你已经发送了一条信息，请耐心等待"

	LimitUserCountTemplate = "Your chat limit has been reached. Invite new users to get 50 chat sessions per new user. your invite link: %s\n\n" +
		"😢您的聊天次数已耗尽，邀请新用户获取聊天次数 30次/新用户. 你的邀请链接: %s"

	RateLimitMessageTemplate = "You can only send %d messages in %d min, please try later. \nRate limiter will disappeared when you join both channel %s and group %s\n\n" +
		"😅 你只能在 %d 分钟内发送 %d 条消息，请稍候再试\n" +
		"当你同时加入频道 %s 和群组 %s 后，将不再限速"

	LimitUserGroupAndChannelTemplate = "Before you join the channel %s and group %s, you can not send private message to me. \n\n" +
		"😅 你必须先加入频道 %s 和群组 %s 才能和我私聊"

	ChatGPTError = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧"
	ChatGPTErrorTemplate = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧 \n\n %s"
	ChatGPTEngineNotOnline = "Chatgpt engine is not ready, please wait a moment. \n\n" +
		"😇ChatGPT 引擎还没有准备好，请稍等一下"
	SendBackMsgFailed = "sendFromChatTask back message failed, please try again later \n\n" +
		"😇返回消息失败，请稍后再试"

	NetworkError = "Network error, please try again later \n\n" +
		"😐网络错误，请稍后再试"

	InternalError = "Internal error, please try again later \n\n" +
		"😐内部错误，请稍后再试"

	ExceedMaxGenerateInviteCodeTimes = "You have exceeded the maximum number of times to generate invite code, please try again later \n\n" +
		"😐你已经超过了生成邀请码的最大次数，请稍后再试"
)
