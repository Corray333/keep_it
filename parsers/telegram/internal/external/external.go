package external

import (
	"bytes"
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

func (e *External) GetTgPhoto(photos []tgbotapi.PhotoSize) (io.Reader, error) {
	photo := photos[len(photos)-1]

	// Get the file URL
	file, err := e.Bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		slog.Error("error getting file: " + err.Error())
		return nil, err
	}

	// Download the file
	fileURL := file.Link(e.Bot.Token)
	resp, err := http.Get(fileURL)
	if err != nil {
		slog.Error("error loading file: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var photoFile bytes.Buffer
	_, err = io.Copy(&photoFile, resp.Body)
	if err != nil {
		slog.Error("error copying file: " + err.Error())
		return nil, err
	}

	return &photoFile, nil
}
