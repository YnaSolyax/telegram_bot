package redis

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type UserRedis struct {
	UserID   int64
	Username string
	Status   int
}

type RedisStorage struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisStorage(logger *zap.Logger) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Could not connect to Redis", zap.Error(err))
	}

	return &RedisStorage{
		client: client,
		logger: logger,
	}
}

func (r *RedisStorage) PutUser(u *UserRedis) error {
	key := strconv.Itoa(int(u.UserID))

	userData, err := json.Marshal(u)
	if err != nil {
		r.logger.Error("Error marshaling user", zap.Int64("userID", u.UserID), zap.Error(err))
		return errors.Wrap(err, "marshal struct")
	}

	err = r.client.Set(ctx, key, userData, 0).Err()
	if err != nil {
		r.logger.Error("Error setting user in Redis", zap.Int64("userID", u.UserID), zap.Error(err))
		return errors.Wrap(err, "set user")
	}

	r.logger.Debug("User saved to Redis", zap.Int64("userID", u.UserID))
	return nil
}

func (r *RedisStorage) GetUser(userID int64) (*UserRedis, error) {
	user := UserRedis{}
	key := strconv.Itoa(int(userID))

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Error checking existence of user in Redis", zap.Int64("userID", userID), zap.Error(err))
		return nil, errors.Wrap(err, "exist user")
	}

	if exists == 0 {
		r.logger.Debug("User does not exist in Redis", zap.Int64("userID", userID))
		return nil, nil
	}

	userData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		r.logger.Error("Error getting user from Redis", zap.Int64("userID", userID), zap.Error(err))
		return nil, errors.Wrap(err, "get user")
	}

	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		r.logger.Error("Error unmarshaling user data", zap.Int64("userID", userID), zap.Error(err))
		return nil, errors.Wrap(err, "unmarshal struct")
	}

	r.logger.Debug("User retrieved from Redis", zap.Int64("userID", userID))
	return &user, nil
}
