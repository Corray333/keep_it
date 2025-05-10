package transport

import (
	"context"
	"log/slog"

	"github.com/Corray333/keep_it/parsers/telegram/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Transport struct {
	*telegram.TelegramClient
	service service
}

type service interface {
	ParseMessage(ctx context.Context, message *tgbotapi.Message) error
}

func New(service service, tgClient *telegram.TelegramClient) *Transport {

	return &Transport{
		tgClient,
		service,
	}
}

func (t *Transport) Run() {
	t.Bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go func() {
			err := t.service.ParseMessage(context.Background(), update.Message)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we failed to keep it😥")
				_, err := t.Bot.Send(msg)
				if err != nil {
					slog.Error("Failed to send error message: ", "error", err)
				}
			}
		}()
	}
}
