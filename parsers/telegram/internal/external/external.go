package external

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/parsers/telegram/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type External struct {
	*telegram.TelegramClient
}

func New(tgClient *telegram.TelegramClient) *External {
	return &External{tgClient}
}

func (e *External) GetTgPhoto(photos []tgbotapi.PhotoSize) ([]byte, error) {
	photo := photos[len(photos)-1]

	// Get the file URL
	file, err := e.Bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		slog.Error("error getting file: ", "error", err)
		return nil, err
	}

	// Download the file
	fileURL := file.Link(e.Bot.Token)
	resp, err := http.Get(fileURL)
	if err != nil {
		slog.Error("error loading file: ", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Rrror reading file", "error", err)
		return nil, err
	}

	return result, nil
}
