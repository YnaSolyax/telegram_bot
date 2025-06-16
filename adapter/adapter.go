package adapter

import (
	"telegram_bot/inmemory"
	"telegram_bot/redis"
	"telegram_bot/storageManager"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CacheAdapter struct {
	cs     *inmemory.InmemoryStorage
	logger *zap.Logger
}

func NewCacheAdapter(cs *inmemory.InmemoryStorage, logger *zap.Logger) *CacheAdapter {
	return &CacheAdapter{
		cs:     cs,
		logger: logger,
	}
}

func (c *CacheAdapter) SetUser(user *storageManager.AdapterUser) error {

	newuser := &inmemory.CacheUser{
		Id:       user.Id,
		Username: user.Username,
		Status:   user.Status,
	}

	err := c.cs.SetUser(newuser)
	if err != nil {
		c.logger.Error("setuser cache adapter error", zap.Error(err))
		return err
	}

	return nil
}

func (c *CacheAdapter) GetUser(id int64) (*storageManager.AdapterUser, error) {
	user, err := c.cs.GetUser(id)

	if err != nil {
		return nil, errors.Wrap(err, "get user cache adapter")
	}

	if user == nil {
		return nil, nil
	}

	newuser := &storageManager.AdapterUser{
		Id:       user.Id,
		Username: user.Username,
		Status:   user.Status,
	}

	return newuser, nil
}

type RedisAdapter struct {
	rs     *redis.RedisStorage
	logger *zap.Logger
}

func NewRedisAdapter(rs *redis.RedisStorage, logger *zap.Logger) *RedisAdapter {
	return &RedisAdapter{
		rs:     rs,
		logger: logger,
	}
}

func (r *RedisAdapter) SetUser(user *storageManager.AdapterUser) error {

	newuser := &redis.UserRedis{
		Id:       user.Id,
		Username: user.Username,
		Status:   user.Status,
	}

	err := r.rs.SetUser(newuser)

	if err != nil {
		r.logger.Error("setuser redis adapter error", zap.Error(err))
		return err
	}

	return nil
}

func (r *RedisAdapter) GetUser(id int64) (*storageManager.AdapterUser, error) {
	user, err := r.rs.GetUser(id)
	if err != nil {
		return nil, errors.Wrap(err, "getuser redis adapter")
	}

	if user == nil {
		return nil, nil
	}

	newuser := &storageManager.AdapterUser{
		Id:       user.Id,
		Username: user.Username,
		Status:   user.Status,
	}

	return newuser, nil
}
