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

type unbanHandler struct {
	Bot    *bot.Bot
	logger *zap.Logger
}

func NewUnBanHandler(b *bot.Bot, logger *zap.Logger) *unbanHandler {
	return &unbanHandler{
		Bot:    b,
		logger: logger,
	}
}

func (ubh *unbanHandler) Handle() {

	f := func(ctx context.Context, b *bot.Bot, update *models.Update) {

		chatID := update.ChannelPost.ReplyToMessage.Chat.ID
		message := update.ChannelPost.ReplyToMessage.Text

		ubh.logger.Debug("Message", zap.Any("details", message))

		re := regexp.MustCompile(`#ID(\d+)`)
		findID := re.FindString(message)

		ok := strings.HasPrefix(findID, "#ID")
		if !ok {
			ubh.logger.Debug("hasprefix unban error", zap.Any("details", findID))
			return
		}

		strId := strings.TrimPrefix(findID, "#ID")
		id, err := strconv.Atoi(strId)

		if err != nil {
			ubh.logger.Debug("Atoi unban error")
			return
		}

		ubh.logger.Debug("", zap.Any("srdId: ", strId), zap.Any("id: ", id), zap.Any("chatd_id:", update.ChannelPost.ReplyToMessage.Chat.ID))

		myb, err := ubh.Bot.UnbanChatMember(ctx, &bot.UnbanChatMemberParams{
			ChatID:       chatID,
			UserID:       int64(id),
			OnlyIfBanned: true,
		})

		ubh.logger.Info("unbanned", zap.Any("id", id), zap.Any("unban?", myb), zap.Any("err?", err))
	}

	ubh.Bot.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return update.ChannelPost != nil &&
				strings.Contains(update.ChannelPost.Text, "/unban")
		},
		f,
	)

}
