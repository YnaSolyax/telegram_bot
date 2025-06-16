package inmemory

import (
	"go.uber.org/zap"
)

type CacheUser struct {
	Id       int64
	Username string
	Status   int
}

type InmemoryStorage struct {
	userMap map[int64]*CacheUser
	logger  *zap.Logger
}

func NewInmemoryStorage(logger *zap.Logger) *InmemoryStorage {
	return &InmemoryStorage{
		userMap: make(map[int64]*CacheUser, 10),
		logger:  logger,
	}
}

func (c *InmemoryStorage) SetUser(u *CacheUser) error {

	c.userMap[u.Id] = &CacheUser{
		Id:       u.Id,
		Username: u.Username,
		Status:   u.Status,
	}

	c.logger.Debug("Cache dump",
		zap.Any("contents", c.userMap),
		zap.Int("total_users", len(c.userMap)))

	return nil
}

func (c *InmemoryStorage) GetUser(id int64) (*CacheUser, error) {

	if c.userMap[id] == nil {
		c.logger.Debug("Cache miss",
			zap.Int64("userID", id),
			zap.Int("cache_size", len(c.userMap)))
		return nil, nil
	}

	c.logger.Debug("Cache dump",
		zap.Any("contents", c.userMap),
		zap.Int("total_users", len(c.userMap)))

	return c.userMap[id], nil
}
