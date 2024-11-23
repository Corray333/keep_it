package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"regexp"

	"github.com/Corray333/keep_it/parsers/vk/internal/entities"
	"github.com/SevereCloud/vksdk/v3/object"
)

type external interface {
	GetTgPhoto(photo *object.PhotosPhoto) (io.Reader, error)
}

type fileManager interface {
	SaveImage(img io.Reader, name string) (string, error)
}

type repository interface {
	NewNote(ctx context.Context, note *entities.NewNoteMessage) error

	SaveNote(ctx context.Context, creationDate, chatID int64, note *entities.Note) error
	GetNotes(ctx context.Context, creationDate, chtID int64) ([]*entities.Note, error)
}

type Service struct {
	external   external
	fileManger fileManager
	repo       repository
}

func New(repo repository, external external, fileManager fileManager) *Service {
	return &Service{
		repo:       repo,
		external:   external,
		fileManger: fileManager,
	}
}

func (s *Service) ParseMessage(ctx context.Context, message *object.MessagesMessage) error {
	// Создаем заметку
	note := &entities.Note{
		Icon:   json.RawMessage(`{"data":"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1em' height='1em' viewBox='0 0 24 24'%3E%3Cpath fill='%2324A1DE' d='M4.26 4.26C3 5.532 3 7.566 3 11.64v.72c0 4.068 0 6.102 1.26 7.38C5.532 21 7.566 21 11.64 21h.72c4.068 0 6.102 0 7.38-1.26C21 18.468 21 16.434 21 12.36v-.72c0-4.068 0-6.102-1.26-7.38C18.468 3 16.434 3 12.36 3h-.72C7.572 3 5.538 3 4.26 4.26m1.776 4.218H8.1c.066 3.432 1.578 4.884 2.778 5.184V8.478h1.938v2.958c1.182-.126 2.43-1.476 2.85-2.964h1.932a5.72 5.72 0 0 1-2.628 3.738a5.92 5.92 0 0 1 3.078 3.756h-2.13c-.456-1.422-1.596-2.526-3.102-2.676v2.676h-.24c-4.104 0-6.444-2.808-6.54-7.488'/%3E%3C/svg%3E"}`),
		Source: "vk",
	}

	// Определяем заголовок и оригинал сообщения
	if message.FwdMessages != nil {
		note.Title = "Forwarded Message"
		// note.Original = fmt.Sprintf("https://vk.com/id%d?w=wall%d_%d", message.FromID, message.PeerID, message.ConversationMessageID)
		note.CopiedAt = int64(message.Date)
	} else {
		note.Title = "Untitled"
		note.CopiedAt = int64(message.Date)
	}

	bottomPhotos := []entities.ImgElement{}

	// Обрабатываем вложения
	for _, attachment := range message.Attachments {
		switch attachment.Type {
		case "photo":
			photo := attachment.Photo
			if len(photo.Sizes) > 0 {
				photoData, err := s.external.GetTgPhoto(&photo)
				if err != nil {
					slog.Error("Error getting photo: " + err.Error())
					return err
				}

				filePath, err := s.fileManger.SaveImage(photoData, "")
				if err != nil {
					slog.Error("Error saving photo: " + err.Error())
					return err
				}
				note.Cover = filePath
			}
		case "wall":
			note.ContentDecoded = append(note.ContentDecoded, entities.TextElement{
				Type:     entities.ElementTypeParagraph,
				RichText: parseRichText(attachment.Wall.Text),
			})
			if note.Original == "" {
				note.Original = fmt.Sprintf("https://vk.com/wall%d_%d", attachment.Wall.OwnerID, attachment.Wall.ID)
			}

			if note.CreatedAt == 0 {
				note.CreatedAt = int64(attachment.Wall.Date)
			}

			for _, wallAttachment := range attachment.Wall.Attachments {
				switch wallAttachment.Type {
				case "photo":
					photo := wallAttachment.Photo
					if len(photo.Sizes) > 0 {
						photoData, err := s.external.GetTgPhoto(&photo)
						if err != nil {
							slog.Error("Error getting photo: " + err.Error())
							return err
						}

						filePath, err := s.fileManger.SaveImage(photoData, "")
						if err != nil {
							slog.Error("Error saving photo: " + err.Error())
							return err
						}
						if note.Cover == "" {
							note.Cover = filePath
						} else {
							bottomPhotos = append(bottomPhotos, entities.ImgElement{
								Type:  entities.ElementTypeImage,
								Src:   filePath,
								Width: 50,
								Align: "center",
							})
						}
					}
				}
			}
		}
	}

	if message.Text != "" {
		note.ContentDecoded = append(note.ContentDecoded, entities.TextElement{
			Type:     entities.ElementTypeParagraph,
			RichText: parseRichText(message.Text),
		})
	}

	for i := range bottomPhotos {
		note.ContentDecoded = append(note.ContentDecoded, bottomPhotos[i])
	}

	contentBytes, err := json.Marshal(note.ContentDecoded)
	if err != nil {
		return fmt.Errorf("failed to marshal caption element: %w", err)
	}
	note.Content = contentBytes

	// Сохраняем заметку
	if err := s.repo.NewNote(ctx, &entities.NewNoteMessage{
		Note:   *note,
		Source: "tg",
		UserID: "377742748",
	}); err != nil {
		return err
	}

	return nil
}

func parseRichText(input string) (result entities.RichText) {
	var meta []entities.Meta
	offsetCorrection := 0

	linkRegex := regexp.MustCompile(`\[(https?://[^\|]+)\|([^\]]+)\]`)
	linkMatches := linkRegex.FindAllStringSubmatchIndex(input, -1)
	for i := range linkMatches {
		if i == 0 && linkMatches[i][0] > 0 {
			result.PlainText += input[offsetCorrection:linkMatches[i][0]]
		}
		link := input[linkMatches[i][2]:linkMatches[i][3]]
		text := input[linkMatches[i][4]:linkMatches[i][5]]
		meta = append(meta, entities.Meta{
			Offset: len([]rune(result.PlainText)) + 1,
			Length: len([]rune(text)),
			Link:   link,
		})

		result.PlainText += text
		if i < len(linkMatches)-1 {
			result.PlainText += input[linkMatches[i][1]:linkMatches[i+1][0]]
		} else {
			result.PlainText += input[linkMatches[i][1]:]
		}
	}

	if len(meta) == 0 {
		result.PlainText = input
	}

	// После обработки метаинформации оставшийся текст является plainText
	result.Meta = meta

	fmt.Printf("Result: %+v\n", result)
	return result
}
