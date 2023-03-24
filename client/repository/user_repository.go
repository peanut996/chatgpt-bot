package repository

import (
	"chatgpt-bot/constant/config"
	botError "chatgpt-bot/constant/error"
	"chatgpt-bot/db"
	"chatgpt-bot/model/persist"
	"chatgpt-bot/utils"
	"errors"
)

var (
	UserTableName = "user"
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

func (u *UserRepository) IsAvailable(userID string) (bool, error) {
	user, err := u.GetByUserID(userID)
	if err != nil || user == nil {
		return false, err
	}
	return user.RemainCount > 0, nil
}

func (u *UserRepository) IsExist(userID string) (bool, error) {
	row := u.db.QueryRow("SELECT remain_count FROM user WHERE user_id = ? LIMIT 1", userID)
	var count int
	err := row.Scan(&count)
	if err != nil && utils.IsNotEmptyRow(err) {
		return false, err
	}
	if utils.IsEmptyRow(err) {
		return false, nil
	}
	return true, nil
}

func (u *UserRepository) generateUniqueInviteCode() (string, error) {
	inviteCode, _ := utils.GenerateInvitationCode(10)
	for i := 0; i < 0xff; i++ {
		user, err := u.GetUserByInviteCode(inviteCode)
		if err != nil {
			return "", err
		}
		if user == nil {
			return inviteCode, nil
		}

		inviteCode, _ = utils.GenerateInvitationCode(10)
	}
	return "", errors.New(botError.ExceedMaxGenerateInviteCodeTimes)
}

func (u *UserRepository) InitUser(userID string, userName string) error {
	inviteCode, err := u.generateUniqueInviteCode()
	if err != nil {
		return err
	}

	_, err = u.db.Exec("INSERT OR IGNORE INTO user (user_id, remain_count, invite_code, user_name) VALUES (?, ?, ?, ?)",
		userID, config.DefaultCount, inviteCode, userName)
	return err
}

func (u *UserRepository) GetByUserID(userID string) (*persist.User, error) {
	user := &persist.User{}
	user.UserID = userID
	row := u.db.QueryRow("SELECT remain_count, invite_code, user_name, is_donate FROM user WHERE user_id = ? LIMIT 1", userID)

	err := row.Scan(&user.RemainCount, &user.InviteCode, &user.UserName, &user.IsDonate)
	if err != nil && utils.IsNotEmptyRow(err) {
		return nil, err
	}
	if utils.IsEmptyRow(err) {
		return nil, nil
	}
	return user, nil
}

func (u *UserRepository) DecreaseCount(userID string) error {
	// check user exist
	exist, err := u.IsExist(userID)
	if err != nil || !exist {
		return err
	}
	_, err = u.db.Exec("UPDATE user SET remain_count = remain_count - 1 WHERE user_id = ?", userID)
	return err
}

func (u *UserRepository) AddCountWhenInviteOther(userID string) error {
	_, err := u.db.Exec("UPDATE user SET remain_count = remain_count + ? WHERE user_id = ?", config.CountWhenInviteOtherUser, userID)
	return err
}

func (u *UserRepository) GetCount(userID string) (int64, error) {

	user, err := u.GetByUserID(userID)
	if err != nil {
		return 0, err
	}
	return user.RemainCount, nil
}

func (u *UserRepository) IsRemainCountMoreThanZero(userID string) (bool, error) {
	count, err := u.GetCount(userID)
	if err != nil {
		return false, err
	}
	return count > 0, nil

}

func (u *UserRepository) UpdateInviteLink(userID, inviteLink string) error {
	_, err := u.db.Exec("UPDATE user SET invite_code = ? WHERE user_id = ?", inviteLink, userID)
	return err
}

func (u *UserRepository) GetUserByInviteCode(inviteCode string) (*persist.User, error) {
	user := &persist.User{}
	user.InviteCode = inviteCode
	row := u.db.QueryRow("SELECT user_id, remain_count, user_name FROM user WHERE invite_code = ? LIMIT 1", inviteCode)

	err := row.Scan(&user.UserID, &user.RemainCount, &user.UserName)
	if err != nil && utils.IsNotEmptyRow(err) {
		return nil, err
	}
	if utils.IsEmptyRow(err) {
		return nil, nil
	}

	return user, nil
}

func (u *UserRepository) GetInviteCodeByUserID(userId string) (string, error) {
	// query user link from db
	user, err := u.GetByUserID(userId)
	if err != nil {
		return "", err
	}
	return user.InviteCode, nil

}

func (u *UserRepository) UpdateCountByUserID(userID string, count string) error {
	num, _ := utils.StringToInt64(count)
	_, err := u.db.Exec("UPDATE user SET remain_count = ? WHERE user_id = ?", num, userID)
	return err

}

func (u *UserRepository) GetAllUserID() ([]string, error) {
	rows, err := u.db.Query("SELECT user_id FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil

}

func (u *UserRepository) Count() (int, error) {
	row := u.db.QueryRow("SELECT count(*) FROM user")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *UserRepository) UpdateUserToVIP(userID string) error {
	_, err := u.db.Exec("UPDATE user SET is_donate = 1 WHERE user_id = ?", userID)
	return err
}

func (u *UserRepository) UpdateUserName(name string, userID string) error {
	_, err := u.db.Exec("UPDATE user SET user_name = ? WHERE user_id = ?", name, userID)
	return err
}
