package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
	"time"
)

func main() {
	adminID := GetEnv("TG_ADMIN_ID")
	token := GetEnv("TG_BOT_TOKEN")

	var telegramBot *bot.Bot
	var ctx context.Context
	var cancel context.CancelFunc
	var notify func(message string) error

	if token != "" {
		telegramBot = TelegramBotAuth(token, adminID)
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		defer telegramBot.Start(ctx)

		notify = func(message string) error {
			_, err := sendMessage(telegramBot, ctx, adminID, message)
			if err != nil {
				return err
			}
			return nil
		}
		err := notify("Bot started.")
		if err != nil {
			panic(err)
		}
	}

	headhunter, err := AuthorizeHeadHunter(GetEnv("HH_USERNAME"), GetEnv("HH_PASSWORD"), notify)
	if err != nil {
		panic(err)
	}

	hoursStr := GetEnv("HH_UPDATE_HOURS")
	hours := StringToInt64(hoursStr)
	initTask(headhunter, hours)
}

func sendMessage(telegramBot *bot.Bot, ctx context.Context, adminID string, text string) (*models.Message, error) {
	return telegramBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: adminID,
		Text:   text,
	})
}

func initTask(headhunter *HeadHunterClient, hours int64) {
	go func() {
		err := headhunter.RaiseResume()
		if err != nil {
			println(err.Error())
		}

		ticker := time.NewTicker(time.Hour * time.Duration(hours))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := headhunter.RaiseResume()
				if err != nil {
					println(err.Error())
				}
			}
		}
	}()
}
