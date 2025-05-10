package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/Corray333/keep_it/parsers/vk/internal/entities"
	"github.com/SevereCloud/vksdk/v3/api"
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
		note.CopiedAt = time.Unix(int64(message.Date), 0)
	} else {
		note.Title = "Untitled"
		note.CopiedAt = time.Unix(int64(message.Date), 0)
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
					slog.Error("Error getting photo: ", "error", err)
					return err
				}

				filePath, err := s.fileManger.SaveImage(photoData, "")
				if err != nil {
					slog.Error("Error saving photo: ", "error", err)
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

			if note.CreatedAt.IsZero() {
				note.CreatedAt = time.Unix(int64(attachment.Wall.Date), 0)
			}

			groupID := strconv.Itoa(-attachment.Wall.FromID)

			vk := api.NewVK(os.Getenv("VK_BOT_TOKEN"))

			groupInfo, err := vk.GroupsGetByID(api.Params{
				"group_id": groupID,
			})
			if err != nil {
				slog.Error("Error getting group info: ", "error", err)
				return err
			}

			if len(groupInfo.Groups) > 0 {
				groupName := groupInfo.Groups[0].Name
				note.Title = "Forwarded from " + groupName
			}

			for _, wallAttachment := range attachment.Wall.Attachments {
				switch wallAttachment.Type {
				case "photo":
					photo := wallAttachment.Photo
					if len(photo.Sizes) > 0 {
						photoData, err := s.external.GetTgPhoto(&photo)
						if err != nil {
							slog.Error("Error getting photo: ", "error", err)
							return err
						}

						filePath, err := s.fileManger.SaveImage(photoData, "")
						if err != nil {
							slog.Error("Error saving photo: ", "error", err)
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
		Source: "vk",
		UserID: strconv.Itoa(message.FromID),
	}); err != nil {
		return err
	}

	return nil
}

func parseRichText(input string) (result entities.RichText) {
	var meta []entities.Meta
	offsetCorrection := 0

	// Регулярка для скрытых гиперссылок
	linkRegex := regexp.MustCompile(`\[(https?://[^\|]+)\|([^\]]+)\]`)
	// Регулярка для явных ссылок
	plainLinkRegex := regexp.MustCompile(`\b(https?://[^\s]+|www\.[^\s]+)\b`)

	linkMatches := linkRegex.FindAllStringSubmatchIndex(input, -1)
	plainLinkMatches := plainLinkRegex.FindAllStringSubmatchIndex(input, -1)

	allMatches := append(linkMatches, plainLinkMatches...)

	// Объединяем результаты поиска ссылок
	sort.Slice(allMatches, func(i, j int) bool {
		return allMatches[i][0] < allMatches[j][0]
	})

	for i := range allMatches {
		match := allMatches[i]

		// Добавляем текст до текущей ссылки, если он есть
		if match[0] > offsetCorrection {
			result.PlainText += input[offsetCorrection:match[0]]
		}

		// Проверяем количество групп в совпадении
		if len(match) == 4 { // Явная ссылка
			link := input[match[0]:match[1]]
			meta = append(meta, entities.Meta{
				Offset: len([]rune(result.PlainText)),
				Length: len([]rune(link)),
				Link:   link,
			})
			result.PlainText += link
		} else if len(match) >= 6 { // Скрытая ссылка
			link := input[match[2]:match[3]]
			text := input[match[4]:match[5]]
			meta = append(meta, entities.Meta{
				Offset: len([]rune(result.PlainText)),
				Length: len([]rune(text)),
				Link:   link,
			})
			result.PlainText += text
		}

		// Устанавливаем новую позицию смещения
		offsetCorrection = match[1]
	}

	// Добавляем оставшийся текст, если он есть
	if offsetCorrection < len(input) {
		result.PlainText += input[offsetCorrection:]
	}

	// Если метаинформации нет, весь текст — plainText
	if len(meta) == 0 {
		result.PlainText = input
	}

	result.Meta = meta

	fmt.Printf("Result: %+v\n", result)
	return result
}
