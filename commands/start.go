package commands

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type StartHandler struct {
	Bot    *bot.Bot
	logger *zap.Logger
}

func NewStartBotHandler(b *bot.Bot, logger *zap.Logger) *StartHandler {
	return &StartHandler{
		Bot:    b,
		logger: logger,
	}
}

func (bh *StartHandler) Handle() {
	f := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		username := update.Message.From.Username
		greeting := fmt.Sprintf("Здравствуйте, @%s !\nНапишите свой вопрос и вам ответят в ближайшее время.", username)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   greeting,
		})

		if err != nil {
			bh.logger.Debug("error send message", zap.Any("details:", err))
			return
		}
	}

	bh.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, f)
}
