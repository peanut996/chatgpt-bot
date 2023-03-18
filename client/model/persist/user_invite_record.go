package persist

import (
	"chatgpt-bot/utils"
	"time"
)

type UserInviteRecord struct {
	UserID string

	InviteUserID string

	InviteTime string
}

func NewUserInviteRecord(userID string, inviteUserID string) *UserInviteRecord {
	return &UserInviteRecord{
		UserID:       userID,
		InviteUserID: inviteUserID,
		InviteTime:   utils.ConvertInt64ToString(time.Now().UnixMilli()),
	}
}
