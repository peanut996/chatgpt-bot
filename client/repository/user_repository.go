package repository

import (
	"chatgpt-bot/db"
	"chatgpt-bot/model/persist"
)

type UserRepository struct {
	db        db.BotDB
	tableName string
}

func NewUserRepository(db db.BotDB) *UserRepository {
	return &UserRepository{
		db:        db,
		tableName: "user",
	}
}

func (u *UserRepository) GetByUserID(userID string) (*persist.User, error) {
	user := &persist.User{}
	user.UserID = userID
	row := u.db.QueryRow("SELECT count, invite_link FROM users WHERE user_id = ? LIMIT 1", userID)

	err := row.Scan(&user.Count, &user.InviteLink)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) IncreaseCount(userID string) error {
	_, err := u.db.Exec("UPDATE user SET count = count + 1 WHERE user_id = ?", userID)
	return err
}

func (u *UserRepository) DecreaseCount(userID string) error {
	_, err := u.db.Exec("UPDATE user SET count = count - 1 WHERE user_id = ?", userID)
	return err
}

func (u *UserRepository) GetCount(userID string) (int64, error) {

	user, err := u.GetByUserID(userID)
	if err != nil {
		return 0, err
	}
	return user.Count, nil
}

func (u *UserRepository) IsRemainCountMoreThanZero(userID string) (bool, error) {
	count, err := u.GetCount(userID)
	if err != nil {
		return false, err
	}
	return count > 0, nil

}

func (u *UserRepository) UpdateInviteLink(userID, inviteLink string) error {
	_, err := u.db.Exec("UPDATE user SET invite_link = ? WHERE user_id = ?", inviteLink, userID)
	return err
}

func (u *UserRepository) FindUserByInviteLink(inviteLink string) (*persist.User, error) {
	user := &persist.User{}
	row := u.db.QueryRow("SELECT user_id, count FROM user WHERE invite_link = ? LIMIT 1", inviteLink)

	err := row.Scan(&user.UserID, &user.Count)
	if err != nil {
		return nil, err
	}
	return user, nil
}
