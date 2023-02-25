package utils

import "strings"

func IsEmpty(s string) bool {
	return s == ""
}

func IsAnyStringEmpty(ss ...string) bool {
	for _, s := range ss {
		if IsEmpty(s) {
			return true
		}
	}
	return false
}
func SplitMessageByMaxSize(msg string, maxSize int) []string {
	var msgs []string
	currentMsg := msg

	if len(currentMsg) <= maxSize {
		msgs = append(msgs, currentMsg)
		return msgs
	}

	for len(currentMsg) > maxSize {
		msgs = append(msgs, currentMsg[:maxSize])
		currentMsg = currentMsg[maxSize:]
	}
	msgs = append(msgs, currentMsg)
	return msgs
}

func ParseBoolString(cmd string) bool {

	cmd = strings.TrimSpace(cmd)

	if strings.Contains(cmd, "false") ||
		strings.Contains(cmd, "off") {
		return false
	}
	return true
}
