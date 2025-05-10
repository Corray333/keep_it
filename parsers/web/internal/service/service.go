package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Corray333/keep_it/parsers/web/internal/entities"
	"github.com/PuerkitoBio/goquery"
)

var NoteWebIcon = json.RawMessage(`{"data": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"14\" height=\"14\" viewBox=\"0 0 14 14\"><g fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1\"><path d=\"M7 13.5a6.5 6.5 0 1 0 0-13a6.5 6.5 0 0 0 0 13M.5 7h13\"/><path d=\"M9.5 7A11.22 11.22 0 0 1 7 13.5A11.22 11.22 0 0 1 4.5 7A11.22 11.22 0 0 1 7 .5A11.22 11.22 0 0 1 9.5 7\"/></g></svg>", "type": "svg"}`)

type fileManager interface {
	SaveFile(ctx context.Context, file []byte, name string) error
}

type repository interface {
	NewNote(ctx context.Context, note *entities.NewNoteMessage) error
	SaveNote(ctx context.Context, creationDate, chatID int64, note *entities.Note) error
	GetNotes(ctx context.Context, creationDate, chtID int64) ([]*entities.Note, error)
}

type Service struct {
	fileManger fileManager
	repo       repository
}

// New создает сервис
func New(repo repository, fileManager fileManager) *Service {
	return &Service{
		repo:       repo,
		fileManger: fileManager,
	}
}

// ProcessHTML обрабатывает HTML, конвертируя в заметку
func (s *Service) ProcessHTML(ctx context.Context, userID int64, document string, url string) error {
	document = "<div>" + document + "</div>"
	note, err := parseHTMLToNote(ctx, document, url)
	if err != nil {
		slog.Error("Error processing page", "error", err)
		return err
	}
	note.CreatorID = userID
	note.Source = "web"
	note.Original = url
	note.Icon = NoteWebIcon

	if err := s.repo.NewNote(ctx, &entities.NewNoteMessage{
		Note:   *note,
		Source: "web",
		UserID: strconv.Itoa(int(userID)),
	}); err != nil {
		return err
	}
	return nil
}

// parseHTMLToNote преобразует HTML и URL в entities.Note
func parseHTMLToNote(ctx context.Context, document string, url string) (*entities.Note, error) {
	if strings.TrimSpace(document) == "" {
		return nil, errors.New("empty document")
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return nil, err
	}

	note := &entities.Note{
		CreatedAt:      time.Now(),
		ContentDecoded: make([]any, 0),
	}

	// Заголовок страницы
	title := strings.TrimSpace(doc.Find("head title").First().Text())
	if title != "" {
		note.Title = title
	}

	doc.Children().Each(func(i int, sel *goquery.Selection) {
		if shouldSkip(sel) {
			return
		}
		switch goquery.NodeName(sel) {
		case "h1", "h2", "h3", "p", "div", "section", "article":
			appendTextWithMeta(note, sel, goquery.NodeName(sel))
			sel.Find("img").Each(func(_ int, img *goquery.Selection) {
				appendImage(note, img)
			})
		case "ul", "ol":
			sel.Find("li").Each(func(_ int, li *goquery.Selection) {
				liText := strings.TrimSpace(li.Text())
				if liText != "" {
					// list items treated as paragraphs
					fake := &goquery.Selection{Nodes: li.Nodes}
					appendTextWithMeta(note, fake, "p")
				}
			})
		case "img":
			appendImage(note, sel)
		default:
			sel.Children().Each(func(_ int, child *goquery.Selection) {
				if !shouldSkip(child) {
					parseHTMLChild(ctx, note, child)
				}
			})
		}
	})

	raw, err := json.Marshal(note.ContentDecoded)
	if err != nil {
		return nil, err
	}
	note.Content = raw
	return note, nil
}

// parseHTMLChild рекурсивно обходит вложенные элементы
func parseHTMLChild(ctx context.Context, note *entities.Note, sel *goquery.Selection) {
	if shouldSkip(sel) {
		return
	}
	switch goquery.NodeName(sel) {
	case "h1", "h2", "h3", "p":
		if goquery.NodeName(sel) == "h1" && note.Title == "" {
			note.Title = strings.TrimSpace(sel.Text())
		}
		appendTextWithMeta(note, sel, goquery.NodeName(sel))
	case "img":
		appendImage(note, sel)
	case "ul", "ol":
		sel.Find("li").Each(func(_ int, li *goquery.Selection) {
			fake := &goquery.Selection{Nodes: li.Nodes}
			appendTextWithMeta(note, fake, "p")
		})
	default:
		if sel.Children().Length() > 0 {
			sel.Children().Each(func(_ int, c *goquery.Selection) {
				parseHTMLChild(ctx, note, c)
			})
		}
	}
}

// appendTextWithMeta добавляет текстовый элемент с метаинформацией (ссылки, жирность, курсив)
func appendTextWithMeta(note *entities.Note, sel *goquery.Selection, tag string) {
	text := strings.TrimSpace(sel.Text())
	if text == "" {
		return
	}
	// базовый элемент типа
	elemType := "p"
	switch tag {
	case "h1", "h2", "h3":
		elemType = tag
	}
	// собрать метаинформацию
	meta := make([]entities.Meta, 0)
	// искать ссылки
	sel.Find("a").Each(func(_ int, a *goquery.Selection) {
		href, ok := a.Attr("href")
		if !ok || strings.TrimSpace(href) == "" {
			return
		}
		anchor := strings.TrimSpace(a.Text())
		if anchor == "" {
			return
		}
		// найти позицию
		idx := strings.Index(text, anchor)
		if idx >= 0 {
			meta = append(meta, entities.Meta{
				Offset: idx,
				Length: len(anchor),
				Link:   href,
			})
		}
	})
	// найти жирный текст
	sel.Find("strong, b").Each(func(_ int, b *goquery.Selection) {
		bold := strings.TrimSpace(b.Text())
		idx := strings.Index(text, bold)
		if idx >= 0 {
			meta = append(meta, entities.Meta{
				Offset: idx,
				Length: len(bold),
				Weight: "bold",
			})
		}
	})
	// найти курсив
	sel.Find("em, i").Each(func(_ int, iel *goquery.Selection) {
		ital := strings.TrimSpace(iel.Text())
		idx := strings.Index(text, ital)
		if idx >= 0 {
			meta = append(meta, entities.Meta{
				Offset: idx,
				Length: len(ital),
				Italic: true,
			})
		}
	})
	// добавить элемент
	note.ContentDecoded = append(note.ContentDecoded, entities.TextElement{
		Type:     entities.ElementType(elemType),
		RichText: entities.RichText{PlainText: text, Meta: meta},
	})
}

// shouldSkip определяет, нужно ли пропустить элемент
func shouldSkip(s *goquery.Selection) bool {
	skipTags := map[string]struct{}{"nav": {}, "aside": {}, "footer": {}, "header": {},
		"button": {}, "form": {}, "input": {}, "select": {},
		"option": {}, "svg": {}, "canvas": {}, "script": {},
	}
	tagName := strings.ToLower(goquery.NodeName(s))
	if _, found := skipTags[tagName]; found {
		return true
	}
	skips := []string{"nav", "menu", "footer", "header", "aside", "btn", "button", "icon", "ads", "share", "sidebar", "navbar", "pagination", "dropdown", "popup"}
	if class, _ := s.Attr("class"); class != "" {
		lower := strings.ToLower(class)
		for _, sk := range skips {
			if strings.Contains(lower, sk) {
				return true
			}
		}
	}
	if id, _ := s.Attr("id"); id != "" {
		lower := strings.ToLower(id)
		for _, sk := range skips {
			if strings.Contains(lower, sk) {
				return true
			}
		}
	}
	if style, _ := s.Attr("style"); style != "" && (strings.Contains(style, "display:none") || strings.Contains(style, "visibility:hidden")) {
		return true
	}
	return false
}

// appendImage добавляет изображение
func appendImage(note *entities.Note, s *goquery.Selection) {
	src, ok := s.Attr("src")
	if !ok || strings.TrimSpace(src) == "" {
		return
	}
	width := 0
	if w, ok := s.Attr("width"); ok {
		if val, err := strconv.Atoi(w); err == nil {
			width = val
		}
	}
	align := ""
	if a, ok := s.Attr("align"); ok {
		align = a
	}
	note.ContentDecoded = append(note.ContentDecoded, entities.ImgElement{
		Type:  entities.ElementTypeImage,
		Src:   src,
		Width: width,
		Align: align,
	})
}
