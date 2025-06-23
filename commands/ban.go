package commands

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type banHandler struct {
	Bot    *bot.Bot
	logger *zap.Logger
}

func NewBanBotHandler(b *bot.Bot, logger *zap.Logger) *banHandler {
	return &banHandler{
		Bot:    b,
		logger: logger,
	}
}

func (bh *banHandler) Handle() {

	f := func(ctx context.Context, b *bot.Bot, update *models.Update) {

		message := update.ChannelPost.ReplyToMessage.Text
		chatdID := update.ChannelPost.ReplyToMessage.Chat.ID

		re := regexp.MustCompile(`#ID(\d+)`)
		regText := re.FindString(message)

		ok := strings.HasPrefix(regText, "#ID")
		if !ok {
			bh.logger.Error("hasprefix error")
			return
		}

		strID := strings.TrimPrefix(regText, "#ID")
		id, err := strconv.Atoi(strID)

		if err != nil {
			bh.logger.Error("Atoi error")
			return
		}

		bh.logger.Debug("", zap.Any("text: ", strID), zap.Any("text: ", id), zap.Any("chatd_id:", update.ChannelPost.ReplyToMessage.Chat.ID))

		myb, err := bh.Bot.BanChatMember(ctx, &bot.BanChatMemberParams{
			ChatID:         chatdID,
			UserID:         int64(id),
			RevokeMessages: true,
		})
		bh.logger.Info("banned", zap.Any("id", id), zap.Any("ban?", myb), zap.Any("err?", err))

	}

	bh.Bot.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update.ChannelPost != nil &&
				strings.Contains(update.ChannelPost.Text, "/ban")
		},
		f,
	)
}
