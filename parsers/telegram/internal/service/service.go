package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/Corray333/keep_it/parsers/telegram/internal/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type external interface {
	GetTgPhoto(photos []tgbotapi.PhotoSize) (io.Reader, error)
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

func (s *Service) ParseMessage(ctx context.Context, message *tgbotapi.Message) error {
	note := &entities.Note{
		Icon:      json.RawMessage(`{"data":"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='1em' height='1em' viewBox='0 0 24 24'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cpath d='m12.593 23.258l-.011.002l-.071.035l-.02.004l-.014-.004l-.071-.035q-.016-.005-.024.005l-.004.01l-.017.428l.005.02l.01.013l.104.074l.015.004l.012-.004l.104-.074l.012-.016l.004-.017l-.017-.427q-.004-.016-.017-.018m.265-.113l-.013.002l-.185.093l-.01.01l-.003.011l.018.43l.005.012l.008.007l.201.093q.019.005.029-.008l.004-.014l-.034-.614q-.005-.018-.02-.022m-.715.002a.02.02 0 0 0-.027.006l-.006.014l-.034.614q.001.018.017.024l.015-.002l.201-.093l.01-.008l.004-.011l.017-.43l-.003-.012l-.01-.01z'/%3E%3Cpath fill='%2324A1DE' d='M19.777 4.43a1.5 1.5 0 0 1 2.062 1.626l-2.268 13.757c-.22 1.327-1.676 2.088-2.893 1.427c-1.018-.553-2.53-1.405-3.89-2.294c-.68-.445-2.763-1.87-2.507-2.884c.22-.867 3.72-4.125 5.72-6.062c.785-.761.427-1.2-.5-.5c-2.302 1.738-5.998 4.381-7.22 5.125c-1.078.656-1.64.768-2.312.656c-1.226-.204-2.363-.52-3.291-.905c-1.254-.52-1.193-2.244-.001-2.746z'/%3E%3C/g%3E%3C/svg%3E"}`),
		CreatedAt: int64(message.Date),
		Source:    "tg",
	}

	if message.ForwardFromChat != nil {
		note.Title = "From " + message.ForwardFromChat.Title
		note.Original = fmt.Sprintf("https://t.me/%s/%d", message.ForwardFromChat.UserName, message.ForwardFromMessageID)
		note.CopiedAt = int64(message.ForwardDate)
	} else {
		note.Title = "Untitled"
		note.CopiedAt = int64(message.Date)
	}

	if message.Photo != nil {
		tgPhoto, err := s.external.GetTgPhoto(message.Photo)
		if err != nil {
			return err
		}

		filePath, err := s.fileManger.SaveImage(tgPhoto, "")
		if err != nil {
			return err
		}

		note.Cover = filePath

	}
	caption := message.Caption

	if caption != "" {
		captionMeta := message.CaptionEntities
		metas := parseEntities(captionMeta)
		captionEl := entities.TextElement{
			Type: entities.ElementTypeParagraph,
			RichText: entities.RichText{
				PlainText: caption,
				Meta:      metas,
			},
		}

		note.ContentDecoded = append(note.ContentDecoded, captionEl)
	}

	if message.Text != "" {
		textMeta := message.Entities
		metas := parseEntities(textMeta)
		textEl := entities.TextElement{
			Type: entities.ElementTypeParagraph,
			RichText: entities.RichText{
				PlainText: message.Text,
				Meta:      metas,
			},
		}

		note.ContentDecoded = append(note.ContentDecoded, textEl)
	}

	if err := s.repo.SaveNote(ctx, int64(message.Date), message.Chat.ID, note); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	notes, err := s.repo.GetNotes(ctx, int64(message.Date), message.Chat.ID)
	if err != nil {
		return err
	}

	if len(notes) == 0 {
		return nil
	}

	finalNote := notes[0]
	for i := 1; i < len(notes); i++ {
		finalNote.ContentDecoded = append(finalNote.ContentDecoded, notes[i].ContentDecoded...)

		if notes[i].Cover != "" {
			if finalNote.Cover == "" {
				finalNote.Cover = notes[i].Cover
			} else {
				finalNote.ContentDecoded = append(finalNote.ContentDecoded, entities.ImgElement{
					Type:  entities.ElementTypeImage,
					Src:   notes[i].Cover,
					Width: 100,
				})
			}
		}

		if notes[i].Original != "" && finalNote.Original == "" {
			finalNote.Original = notes[i].Original
			finalNote.Title = notes[i].Title
		}
	}

	contentBytes, err := json.Marshal(finalNote.ContentDecoded)
	if err != nil {
		return fmt.Errorf("failed to marshal caption element: %w", err)
	}
	finalNote.Content = contentBytes

	newNoteMsg := &entities.NewNoteMessage{
		Note:   *finalNote,
		Source: "tg",
		UserID: strconv.Itoa(int(message.From.ID)),
	}

	if err = s.repo.NewNote(ctx, newNoteMsg); err != nil {
		return err
	}

	return nil
}

func parseEntities(msgEntities []tgbotapi.MessageEntity) []entities.Meta {
	// Сортируем сущности по началу (Offset) для корректного объединения
	sort.Slice(msgEntities, func(i, j int) bool {
		return msgEntities[i].Offset < msgEntities[j].Offset
	})

	var metas []entities.Meta

	for _, entity := range msgEntities {
		if entity.Type == "mention" || entity.Type == "hashtag" || entity.Type == "bot_command" || entity.Type == "custom_emoji" {
			continue
		}

		end := entity.Offset + entity.Length

		// Обрабатываем различные типы сущностей
		meta := entities.Meta{
			Offset: entity.Offset,
			Length: entity.Length,
		}

		switch entity.Type {
		case "bold":
			meta.Weight = "bold"
		case "italic":
			meta.Italic = true
		case "underline":
			meta.Underline = true
		case "strikethrough":
			meta.Strikethrough = true
		}

		// Если сущность имеет URL, добавляем его в Meta
		if entity.URL != "" {
			meta.Link = entity.URL
		}

		// Объединение с предыдущими метками, если они пересекаются
		merged := false
		for i := range metas {
			if metas[i].Offset <= meta.Offset && metas[i].Offset+metas[i].Length >= end {
				metas[i] = mergeMeta(metas[i], meta)
				merged = true
				break
			}
		}

		// Добавляем новую метку, если не удалось объединить
		if !merged {
			metas = append(metas, meta)
		}
	}

	return metas
}

// mergeMeta объединяет два объекта Meta, если они пересекаются
func mergeMeta(a, b entities.Meta) entities.Meta {
	if b.Weight == "bold" {
		a.Weight = "bold"
	}
	if b.Italic {
		a.Italic = true
	}
	if b.Underline {
		a.Underline = true
	}
	if b.Strikethrough {
		a.Strikethrough = true
	}
	if b.Link != "" && a.Link == "" {
		a.Link = b.Link
	}
	return a
}
