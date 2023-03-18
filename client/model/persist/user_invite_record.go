package persist

import (
	"chatgpt-bot/utils"
	"time"
)

type UserInviteRecord struct {
	UserID string

	UserInviteID string

	InviteTime string
}

func NewUserInviteRecord(userID string, userInviteID string, inviteTime string) *UserInviteRecord {
	return &UserInviteRecord{
		UserID:       userID,
		UserInviteID: userInviteID,
		InviteTime:   utils.ConvertInt64ToString(time.Now().UnixMilli()),
	}
}
