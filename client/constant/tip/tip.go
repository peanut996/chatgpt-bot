package tip

var (
	BotPingTip    = "pong"
	UnknownCmdTip = "Unknown command, please send /start to start a chat \n\n" +
		"🔥未知命令，请发送 /start 来开始聊天"
	BotStartTip = "Just send me a sentence and I will reply you. \n" +
		"You can also add me to your own group\n" +
		"Bot default use gpt-3.5 model, if you want to use gpt-4 model, please use `/gpt4` command, like 【/gpt4 how is weather today?】 \n\n" +
		"😊请在这条消息下回复你的问题，我会回复你的 \n" +
		"🔥你也可以私聊我或者把我加到你的群组聊天 \n" +
		"🤖默认使用gpt-3.5模型，gpt-4模型对话请使用「 /gpt4 」命令加上你的问题，如【/gpt4 今天天气怎么样?】"

	NotAdminTip = "You are not admin, can not use this command \n\n" +
		"😢你不是管理员，无法执行此操作"

	InviteSuccessTip = "Invite success, you can use /invite to get your invite link \n\n" +
		"😊邀请成功，你可以使用 /invite 来获取你的邀请链接"

	InviteTipTemplate = "You can invite new users to get 30 chat sessions per new user. your invite link: %s\n\n" +
		"😊你可以邀请新用户获取聊天次数 30次/新用户. 你的邀请链接: %s"

	QueryUserInfoTemplate = "💁账号(Account): %s\n\n" +
		"🏆剩余次数(RemainCount): %d\n" +
		"🎭邀请人数(InviteUsers): %d\n" +
		"🔗邀请链接(InviteLink): %s\n\n" +
		"🔮小提示：邀请1人获得30次聊天次数(Invite 1 user to get 30 chat count)"

	GPT4LackTextTip = "Please add your question after `/gpt4` \n\n" +
		"😊请在 `/gpt4`命令后面加上你的问题"
)
