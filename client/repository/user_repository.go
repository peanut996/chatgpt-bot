package repository

import (
	"chatgpt-bot/db"
	"chatgpt-bot/model/persist"
)

var (
	UserTableName = "user"

	DefaultCount = 30

	CountWhenInviteOtherUser = 30
)

type UserRepository struct {
	db        db.BotDB
	tableName string
}

func NewUserRepository(db db.BotDB) *UserRepository {
	return &UserRepository{
		db:        db,
		tableName: UserTableName,
	}
}

func (u *UserRepository) IsAvaliable(userID string) (bool, error) {
	exist, err := u.IsExist(userID)
	if err != nil {
		return false, err
	}
	if !exist {
		err = u.InitUser(userID)
		if err != nil {
			return false, err
		}
	}
	return u.IsRemainCountMoreThanZero(userID)
}

func (u *UserRepository) IsExist(userID string) (bool, error) {
	row := u.db.QueryRow("SELECT count(*) FROM user WHERE user_id = ? LIMIT 1", userID)
	var count int
	err := row.Scan(&count)
	if err == nil {
		return true, nil
	}
	return false, err
}

func (u *UserRepository) InitUser(userID string) error {
	_, err := u.db.Exec("INSERT OR IGNORE INTO user (user_id, count) VALUES (?, ?)", userID, DefaultCount)
	return err
}

func (u *UserRepository) GetByUserID(userID string) (*persist.User, error) {
	user := &persist.User{}
	user.UserID = userID
	row := u.db.QueryRow("SELECT count, invite_link FROM user WHERE user_id = ? LIMIT 1", userID)

	err := row.Scan(&user.Count, &user.InviteLink)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) DecreaseCount(userID string) error {
	// check user exist
	exist, err := u.IsExist(userID)
	if err != nil {
		return err
	}
	if !exist {
		err = u.InitUser(userID)
		if err != nil {
			return err
		}
	}
	_, err = u.db.Exec("UPDATE user SET count = count - 1 WHERE user_id = ?", userID)
	return err
}

func (u *UserRepository) AddCountWhenInviteOther(userID string) error {
	_, err := u.db.Exec("UPDATE user SET count = count + ? WHERE user_id = ?", CountWhenInviteOtherUser, userID)
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

func (u *UserRepository) GetUserInviteLink(userId string) (string, error) {
	// query user link from db
	user, err := u.GetByUserID(userId)
	if err != nil {
		return "", err
	}
	return user.InviteLink, nil

}
