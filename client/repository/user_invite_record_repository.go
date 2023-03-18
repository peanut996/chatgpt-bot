package repository

import (
	"chatgpt-bot/db"
	"chatgpt-bot/model/persist"
	"chatgpt-bot/utils"
)

var (
	UserInviteRecordTableName = "user_invite_record"
)

type UserInviteRecordRepository struct {
	db        db.BotDB
	tableName string
}

func NewUserInviteRecordRepository(db db.BotDB) *UserInviteRecordRepository {
	return &UserInviteRecordRepository{
		db:        db,
		tableName: UserInviteRecordTableName,
	}
}

func (r *UserInviteRecordRepository) Insert(record *persist.UserInviteRecord) error {
	raw := "INSERT INTO user_invite_record (user_id, invite_user_id, invite_time) VALUES (?, ?, ?)"
	_, err := r.db.Exec(raw, record.UserID, record.InviteUserID, record.InviteTime)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserInviteRecordRepository) CountByUserID(userID string) (int64, error) {
	raw := "SELECT COUNT(*) FROM user_invite_record WHERE user_id = ?"
	var count int64
	err := r.db.QueryRow(raw, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserInviteRecordRepository) GetByInviteUserID(inviteUserID string) (*persist.UserInviteRecord, error) {
	record := &persist.UserInviteRecord{}
	record.InviteUserID = inviteUserID
	row := r.db.QueryRow("SELECT user_id, invite_time FROM user_invite_record WHERE invite_user_id = ? LIMIT 1", inviteUserID)

	err := row.Scan(&record.UserID, &record.InviteTime)
	if err != nil && utils.IsNotEmptyRow(err) {
		return nil, err
	}
	if utils.IsEmptyRow(err) {
		return nil, nil
	}
	return record, nil
}
