package cmd

var (
	START   = "start"
	PING    = "ping"
	CHATGPT = "chatgpt"

	LIMITER = "limiter"

	PPROF = "pprof"

	INVITE = "invite"

	COUNT = "count"

	QUERY = "query"

	GPT4 = "gpt4"

	DONATE = "donate"

	_ = "cmd"
)

func IsBotCmd(cmd string) bool {
	switch cmd {
	case START, PING, CHATGPT, LIMITER, PPROF, INVITE, COUNT, QUERY, GPT4, DONATE:
		return true
	default:
		return false
	}
}
