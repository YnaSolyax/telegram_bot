package commands

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type banHandler struct {
	Bot    *bot.Bot
	ChatID string
}

func NewBanBotHandler(b *bot.Bot, chatID string) *banHandler {
	return &banHandler{
		Bot:    b,
		ChatID: chatID,
	}
}

func (bb *banHandler) Handle() {

	f := func(ctx context.Context, b *bot.Bot, update *models.Update) {

		userToBan := update.Message.ReplyToMessage.From
		userID := userToBan.ID

		if update.Message.ReplyToMessage != nil {

			b.BanChatMember(ctx, &bot.BanChatMemberParams{
				ChatID:         bb.ChatID,
				UserID:         userID,
				RevokeMessages: true,
			})

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: bb.ChatID,
				Text:   "Пользователь @%s забанен",
			})
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: bb.ChatID,
				Text:   "Вы должны ответить на сообщение пользователя, чтобы забанить его.",
			})
		}
	}
	bb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, f)
}
