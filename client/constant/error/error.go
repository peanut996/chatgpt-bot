package error

var (
	MissingRequiredConfig = "missing required config"
	EmptyMessage          = ""

	OnlyOneChatAtATime = "you are chatting with me, please wait for a while. \n\n" +
		"😅你已经发送了一条信息，请耐心等待"
	LimitUserCountTemplate = "Your chat limit has been reached. Invite one new user to get more %d times. your invite link: %s\n\n" +
		"😢您的聊天次数已耗尽，邀请新用户后可获得%d聊天次数. 你的邀请链接: %s"

	RateLimitMessageTemplate = "You are chatting with me too frequently, can only send %d messages in %d min, remain %d seconds.\n" +
		"Tip: you can invite new user to avoid the rate limiter.\n\n" +
		"😅你聊天太频繁了, 只能在 %d 分钟内发送 %d 条消息，还剩 %d 秒\n\n" +
		"🔮小提示：成为捐赠用户或邀请%d个新用户可解除消息速率限制, 使用 「 /invite 」 获取邀请链接"

	LimitUserGroupAndChannel = "Before you join the channel and group, you can not send message to me. \n\n" +
		"😅 你必须先加入频道和群组才能和我聊天"

	ChatGPTError = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧"
	ChatGPTErrorTemplate = "ChatGPT return error, try later again \n\n" +
		"😇出错了, 稍后重试下吧 \n\n%s"
	ChatGPTEngineNotOnline = "Chatgpt engine is not ready, please wait a moment. \n\n" +
		"😇ChatGPT 引擎还没有准备好，请稍等一下"

	NetworkError = "Network error, please try again later \n\n" +
		"😐网络错误，请稍后再试"

	InternalError = "Internal error, please try again later \n\n" +
		"😐内部错误，请稍后再试"

	ExceedMaxGenerateInviteCodeTimes = "You have exceeded the maximum number of times to generate invite code, please try again later \n\n" +
		"😐你已经超过了生成邀请码的最大次数，请稍后再试"

	OnlyAllowInPrivate = "Only allow in private chat \n\n" +
		"😐只允许在私聊中使用"

	OnlyAllowInGroup = "Only allow in group \n\n" +
		"😐只允许在群组中使用"

	OnlyAllowInChannel = "Only allow in channel \n\n" +
		"😐只允许在频道中使用"

	CannotGetAccessCodeTemplate = "Cannot get access code, become a donor or invite %d new users to obtain an access code. Use `/invite` to get an invitation link. \n\n" +
		"😐无法获取访问码，成为捐赠用户或邀请%d个新用户可获取访问码，使用 「 /invite 」 获取邀请链接"
)
