package error

var (
	MissingRequiredConfig = "missing required config"
	EmptyMessage          = ""

	OnlyOneChatAtATime = "you are chatting with me, please wait for a while. \n\n" +
		"😅你已经发送了一条信息，请耐心等待"
	LimitUserCountTemplate = "Your chat limit has been reached. Invite one new user to get more %d times. your invite link: %s\n\n" +
		"😢您的聊天次数已耗尽，邀请新用户后可获得%d聊天次数. 你的邀请链接: %s"

	RateLimitMessageTemplate = "You are chatting with me too frequently, can only send %d messages in %d min, remain %d seconds. \n\n" +
		"😅你聊天太频繁了, 只能在 %d 分钟内发送 %d 条消息，还剩 %d 秒\n"

	LimitUserGroupAndChannelTemplate = "Before you join the channel %s and group %s, you can not send private message to me. \n\n" +
		"😅 你必须先加入频道 %s 和群组 %s 才能和我私聊"

	ChatGPTError = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧"
	ChatGPTErrorTemplate = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧 \n\n%s"
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
