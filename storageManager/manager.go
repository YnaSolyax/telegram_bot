package storageManager

import (
	"fmt"
	"telegram_bot/storage"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type AdapterUser struct {
	Id       int64
	Username string
	Status   int
}

type StorageMethods interface {
	SetUser(*AdapterUser) error
	GetUser(id int64) (*AdapterUser, error)
}

type StorageUser struct {
	logger  *zap.Logger
	manager *storage.DBManager
	sm      StorageMethods
}

func Manager(manager *storage.DBManager, sm StorageMethods, logger *zap.Logger) *StorageUser {
	return &StorageUser{
		logger:  logger,
		manager: manager,
		sm:      sm,
	}
}

func (u *StorageUser) SetUser(userID int64, username string, status int) error {
	u.logger.Debug("Inserting user", zap.Int64("userID", userID), zap.String("username", username))

	user, err := u.GetUser(userID)
	if err != nil {
		return errors.Wrap(err, "get user")
	}

	if user != nil {
		return nil
	}

	if err := u.AddUserToDB(userID, username, status); err != nil {
		return errors.Wrap(err, "add user to db")
	}

	if err := u.sm.SetUser(&AdapterUser{
		Id:       userID,
		Username: username,
		Status:   status,
	}); err != nil {
		return errors.Wrap(err, "put user in redis")
	}

	return nil
}

func (u *StorageUser) GetUser(userID int64) (*AdapterUser, error) {
	u.logger.Debug("Checking user", zap.Int64("userID", userID))

	user, err := u.sm.GetUser(userID)
	if err != nil {
		u.logger.Warn("User not found in cache, checking DB",
			zap.Int64("userID", userID),
			zap.String("source", "cache"))
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	u.logger.Debug("Get user from DB")

	dbUser, err := u.manager.GetUser(userID)
	if err != nil {
		u.logger.Error("Error getting user from DB", zap.Int64("userID", userID), zap.Error(err))
		return nil, fmt.Errorf("ERROR! manager.go:82")
	}

	if dbUser == nil {
		u.logger.Debug("User not found in DB", zap.Int64("userID", userID))
		return nil, nil
	}

	err = u.sm.SetUser(&AdapterUser{
		Id:       dbUser.Id,
		Username: dbUser.Username,
		Status:   dbUser.Status,
	})
	if err != nil {
		u.logger.Error("Error saving user to Redis", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}

	u.logger.Debug("User saved to Storage from DB", zap.Int64("userID", userID))

	return user, nil
}
