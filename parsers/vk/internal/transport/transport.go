package transport

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Corray333/keep_it/parsers/vk/internal/vk"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/object"
)

type Transport struct {
	*vk.VKClient
	service service
}

type service interface {
	ParseMessage(ctx context.Context, message *object.MessagesMessage) error
}

func New(service service, tgClient *vk.VKClient) *Transport {

	return &Transport{
		tgClient,
		service,
	}
}

func (t *Transport) Run() {
	t.Bot.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		message := obj.Message
		if err := t.service.ParseMessage(ctx, &message); err != nil {
			slog.Error("Error parsing message: " + err.Error())
		}
	})

	// Запускаем Long Poll
	fmt.Println("Bot started")
	if err := t.Bot.Run(); err != nil {
		slog.Error("Error running bot: " + err.Error())
	}
}
