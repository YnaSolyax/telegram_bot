package commands

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type StartHandler struct {
	Bot       *bot.Bot
	ChannelID string
}

func NewStartBotHandler(b *bot.Bot, channelID string) *StartHandler {
	return &StartHandler{
		Bot:       b,
		ChannelID: channelID,
	}
}

func (bh *StartHandler) Handle() error {
	f := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		username := update.Message.From.Username
		greeting := fmt.Sprintf("Здравствуйте, @%s !\nНапишите свой вопрос и вам ответят в ближайшее время.", username)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: bh.ChannelID,
			Text:   greeting,
		})

		if err != nil {
			return
		}
	}

	bh.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, f)
	return nil
}
