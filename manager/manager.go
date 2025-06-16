package manager

import (
	"context"
	"os"
	"os/signal"
	"telegram_bot/adapter"
	"telegram_bot/commands"
	"telegram_bot/handler"
	"telegram_bot/inmemory"
	"telegram_bot/middleware"
	"telegram_bot/storage"
	"telegram_bot/storageManager"

	"github.com/go-telegram/bot"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	channelID = "YOUR_CHANAL"
	token     = "YOUR_API_TOKEN"
)

func Manager() error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Error("Error create logger")
		return errors.Wrap(err, "logger error")
	}

	db, err := storage.NewDB()
	if err != nil {
		logger.Error("Error create database")
		return errors.Wrap(err, "db error")
	}

	newDbManager := storage.NewDBManager(db)

	//redis := redis.NewRedisStorage(logger)
	//redisAdapter := adapter.NewRedisAdapter(redis, logger)

	cache := inmemory.NewInmemoryStorage(logger)
	cacheAdapter := adapter.NewCacheAdapter(cache, logger)

	stManager := storageManager.Manager(newDbManager, cacheAdapter, logger)
	userMiddleware := middleware.NewUser(logger, stManager)

	opts := []bot.Option{
		bot.WithMiddlewares(userMiddleware.Handler),
		bot.WithDefaultHandler(handler.Handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		logger.Error("Error get bot token...")
		return errors.Wrap(err, "bot token error")
	}

	commands.NewStartBotHandler(b, channelID).Handle()
	commands.NewBanBotHandler(b, channelID).Handle()
	b.Start(ctx)

	return nil

}
