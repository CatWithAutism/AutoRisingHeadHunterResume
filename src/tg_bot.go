package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"strconv"
)

func GetAdminHandler(adminId string) func(next bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message != nil && update.Message.Text != "" {
				if update.Message.From.ID == StringToInt64(adminId) {
					next(ctx, b, update)
				} else {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "You are not admin.",
					})
				}
			}
		}
	}
}

func TelegramBotAuth(token string, adminId string) *bot.Bot {
	opts := []bot.Option{
		bot.WithMiddlewares(GetAdminHandler(adminId)),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}
	return b
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
