package middleware

import (
	"context"
	"telegram_bot/storageManager"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type User struct {
	logger  *zap.Logger
	storage *storageManager.StorageUser
}

func NewUser(logger *zap.Logger, storage *storageManager.StorageUser) *User {
	return &User{
		logger:  logger,
		storage: storage,
	}
}

func (u *User) Handler(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		userID := update.Message.From.ID
		username := update.Message.From.Username
		status := 0

		err := u.storage.SetUser(userID, username, status)
		if err != nil {
			u.logger.Error("User operation failed",
				zap.Int64("user_id", userID),
				zap.String("username", username),
				zap.String("error", err.Error()))
		}

		next(ctx, bot, update)
	}
}
