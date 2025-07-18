package handler

import (
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var channelID = "@testchannerll"

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.ChannelPost == nil {
		username := update.Message.From.Username
		text := update.Message.Text
		messageID := strconv.FormatInt(update.Message.From.ID, 10)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: channelID,
			Text:   "Сообщение от @" + username + "(#ID" + messageID + "):\n" + text,
		})
	}
}
