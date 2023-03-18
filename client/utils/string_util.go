package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
)

func IsEmpty(s string) bool {
	return s == ""
}

func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
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

func GenerateInvitationCode(size int) (string, error) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, size)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[index.Int64()]
	}
	return string(result), nil
}

func ConvertUserID(userID int64) string {
	return strconv.FormatInt(userID, 10)
}
