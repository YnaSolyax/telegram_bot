package main

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

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger, _ := zap.NewDevelopment()

	db, err := storage.NewDB()
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
		panic(err)
	}

	commands.NewStartBotHandler(b, channelID).Handle()
	commands.NewBanBotHandler(b, channelID).Handle()
	b.Start(ctx)
}
