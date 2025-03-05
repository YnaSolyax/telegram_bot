package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func handleMessageFromChannel(ctx context.Context, b *bot.Bot, update *models.Update) {

	text := update.Message.Text

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: channelID,
		Text:   text,
	})

}
