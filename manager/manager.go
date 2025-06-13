package manager

import (
	"context"
	"os"
	"os/signal"
	"telegram_bot/commands"
	"telegram_bot/handler"
	"telegram_bot/middleware"
	"telegram_bot/redis"
	"telegram_bot/storage"
	"telegram_bot/storageManager"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

var channelID = "Your_chaanel"

func Manager() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Error("Error create logger")
	}

	db, err := storage.NewDB()
	if err != nil {
		logger.Error("Error create database")
	}

	nm := storage.NewDBManager(db)
	redis := redis.NewRedisStorage(logger)
	su := storageManager.Manager(nm, redis, logger)
	userMiddleware := middleware.NewUser(logger, su)

	opts := []bot.Option{
		bot.WithMiddlewares(userMiddleware.Handler),
		bot.WithDefaultHandler(handler.Handler),
	}

	b, err := bot.New("token", opts...)
	if err != nil {
		logger.Error("Error get bot token...")
	}

	commands.NewStartBotHandler(b, channelID).Handle()
	commands.NewBanBotHandler(b, channelID).Handle()
	b.Start(ctx)

}
