package external

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/Corray333/keep_it/parsers/vk/internal/vk"
	"github.com/SevereCloud/vksdk/v3/object"
)

type External struct {
	*vk.VKClient
}

func New(tgClient *vk.VKClient) *External {
	return &External{tgClient}
}

func (e *External) GetTgPhoto(photo *object.PhotosPhoto) (io.Reader, error) {

	// Get the file URL
	fileURL := photo.MaxSize().URL

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
