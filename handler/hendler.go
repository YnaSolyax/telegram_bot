package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var channelID = "@testchannerll"

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		handleMessageFromChannel(ctx, b, update)
	}

	if update.Message.From != nil {
		handleMessageFromUser(ctx, b, update)
	}

}
