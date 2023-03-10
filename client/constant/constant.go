package constant

var ChatGPTTimeoutSeconds = 360

var (
	BotPingTip    = "pong"
	UnknownCmdTip = "Unknown command, please send /start to start a chat \n\n" +
		"ð¥æªç¥å½ä»¤ï¼è¯·åé /start æ¥å¼å§èå¤©"
	BotStartTip = "Hi, I'm ChatGPT bot. I can chat with you. Just send me a sentence and I will reply you. \nYou can also add me to your own group\n\n" +
		"ðè¯·å¨è¿æ¡æ¶æ¯ä¸åå¤ä½ çé®é¢ï¼æä¼åå¤ä½ ç \n\nð¥ä½ ä¹å¯ä»¥ç§èææèææå å°ä½ çç¾¤ç»èå¤©"

	NotAdminTip = "You are not admin, can not use this command \n\n" +
		"ð¢ä½ ä¸æ¯ç®¡çåï¼æ æ³æ§è¡æ­¤æä½"
)

var (
	OnlyOneChatAtATime = "you are chatting with me, please wait for a while. \n\n" +
		"ðä½ å·²ç»åéäºä¸æ¡ä¿¡æ¯ï¼è¯·èå¿ç­å¾"

	LimitUserMessageTemplate = "You should join channel %s and group %s, then you can talk to me \n\n" +
		"ð¢ä½ éè¦å å¥é¢é %s åç¾¤ç» %sï¼ç¶åæè½åæäº¤è°"

	RateLimitMessageTemplate = "You can only send %d messages in %d min, please try later. \nRate limiter will disappeared when you join both channel %s and group %s\n\n" +
		"ð ä½ åªè½å¨ %d åéååé %d æ¡æ¶æ¯ï¼è¯·ç¨ååè¯\n" +
		"å½ä½ åæ¶å å¥é¢é %s åç¾¤ç» %s åï¼å°ä¸åéé"

	ChatGPTError = "ChatGPT return error, try later again \n\n" +
		"ðåºéäº, ç¨åéè¯ä¸å§"
	ChatGPTErrorTemplate = "ChatGPT return error, try later again \n\n" +
		"ðåºéäº, ç¨åéè¯ä¸å§ \n\n %s"
	ChatGPTEngineNotOnline = "Chatgpt engine is not ready, please wait a moment. \n\n" +
		"ðChatGPT å¼æè¿æ²¡æåå¤å¥½ï¼è¯·ç¨ç­ä¸ä¸"
	SendBackMsgFailed = "Send back message failed, please try again later \n\n" +
		"ðè¿åæ¶æ¯å¤±è´¥ï¼è¯·ç¨ååè¯"

	NetworkError = "Network error, please try again later \n\n" +
		"ðç½ç»éè¯¯ï¼è¯·ç¨ååè¯"

	InternalError = "Internal error, please try again later \n\n" +
		"ðåé¨éè¯¯ï¼è¯·ç¨ååè¯"
)
