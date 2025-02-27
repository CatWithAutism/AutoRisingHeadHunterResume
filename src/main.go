package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	httpClient, hunter, _ := AuthorizeHeadHunter(GetEnv("HH_USERNAME"), GetEnv("HH_PASSWORD"))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	adminID := GetEnv("TG_ADMIN_ID")
	token := GetEnv("TG_BOT_TOKEN")

	telegramBot := TelegramBotAuth(token, adminID)
	sendMessage(telegramBot, ctx, adminID, "Bot started.")
	initTask(httpClient, hunter, func() { sendMessage(telegramBot, ctx, adminID, "Resume updated.") })
	telegramBot.Start(ctx)

}

func sendMessage(telegramBot *bot.Bot, ctx context.Context, adminID string, text string) (*models.Message, error) {
	return telegramBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: adminID,
		Text:   text,
	})
}

func initTask(httpClient *http.Client, hunter *HeadHunterCookies, notify func()) {
	go func() {
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				RaiseResume(httpClient, hunter)
				notify()
			}
		}
	}()
}
