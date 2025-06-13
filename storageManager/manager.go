package storageManager

import (
	"fmt"
	"telegram_bot/redis"
	"telegram_bot/storage"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type StorageUser struct {
	logger  *zap.Logger
	manager *storage.DBManager
	redis   *redis.RedisStorage
}

func Manager(manager *storage.DBManager, redisStorage *redis.RedisStorage, logger *zap.Logger) *StorageUser {
	return &StorageUser{
		logger:  logger,
		manager: manager,
		redis:   redisStorage,
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

	if err := u.redis.SetUser(&redis.UserRedis{
		UserID:   userID,
		Username: username,
		Status:   status,
	}); err != nil {
		return errors.Wrap(err, "put user in redis")
	}

	return nil
}

func (u *StorageUser) GetUser(userID int64) (*redis.UserRedis, error) {
	u.logger.Debug("Checking user", zap.Int64("userID", userID))

	user, err := u.redis.GetUser(userID)
	if err != nil {
		u.logger.Error("Error getting user from Redis", zap.Int64("userID", userID), zap.Error(err))
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

	err = u.redis.SetUser((*redis.UserRedis)(dbUser))
	if err != nil {
		u.logger.Error("Error saving user to Redis", zap.Int64("userID", userID), zap.Error(err))
		return nil, err
	}

	u.logger.Debug("User saved to Redis from DB", zap.Int64("userID", userID))

	return user, nil
}
