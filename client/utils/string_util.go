package utils

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

func PanicIfAnyStringEmpty(ss ...string) {
	if IsAnyStringEmpty(ss...) {
		panic("any string is empty")
	}
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
