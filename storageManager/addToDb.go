package storageManager

import "telegram_bot/storage"

func (u *StorageUser) AddUserToDB(userID int64, username string, status int) error {

	user, err := u.manager.GetUser(userID)
	if err != nil {
		return err
	}
	if user != nil {
		return err
	}

	user = &storage.DBUser{
		UserID:   userID,
		Username: username,
		Status:   status,
	}

	u.manager.SetUser(user)

	return nil
}
