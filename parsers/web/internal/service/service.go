package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Corray333/keep_it/parsers/web/internal/entities"
	"github.com/PuerkitoBio/goquery"
)

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

func New(repo repository, fileManager fileManager) *Service {
	return &Service{
		repo:       repo,
		fileManger: fileManager,
	}
}

func isContentBlock(s *goquery.Selection) bool {
	// Список классов, указывающих на контент
	contentClasses := []string{"article", "Article", "blog", "content", "post", "main-content", "story"}
	// Список классов и тегов, указывающих на служебные блоки
	nonContentClasses := []string{"nav", "navbar", "footer", "sidebar", "menu", "button", "widget", "ad", "advertisement"}
	nonContentTags := []string{"nav", "footer", "aside", "header"}

	// Проверяем тег <article>
	if goquery.NodeName(s) == "article" {
		return true
	}

	// Проверяем классы
	classStr, exists := s.Attr("class")
	if exists {
		classes := strings.Fields(classStr)
		for _, class := range classes {
			for _, contentClass := range contentClasses {
				if strings.Contains(class, contentClass) {
					return true
				}
			}
			for _, nonContentClass := range nonContentClasses {
				if strings.Contains(class, nonContentClass) {
					return false
				}
			}
		}
	}

	// Проверяем тег на служебность
	tagName := goquery.NodeName(s)
	for _, nonContentTag := range nonContentTags {
		if tagName == nonContentTag {
			return false
		}
	}

	// Эвристика: если блок содержит много текста (более 100 символов), считаем его контентным
	text := strings.TrimSpace(s.Text())
	if len(text) > 100 {
		return true
	}

	return false
}

func (s *Service) ProcessHTML(ctx context.Context, userID int64, document string, url string) error {
	note, err := parseHTMLToNote(document, url)
	if err != nil {
		return err
	}

	fmt.Println(note)
	return nil

	if err := s.repo.NewNote(ctx, &entities.NewNoteMessage{
		Note:   *note,
		Source: "web",
		UserID: strconv.Itoa(int(userID)),
	}); err != nil {
		return err
	}

	return nil

}

// Функция парсинга HTML в структуру Note
func parseHTMLToNote(document string, url string) (*entities.Note, error) {
	// Создаем goquery документ из строки HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Массив для элементов контента
	var contentElements []interface{}

	// Находим все потенциальные контентные блоки
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		if isContentBlock(s) {
			// Парсим только внутри этого блока
			s.Find("h1, h2, h3, p").Each(func(i int, sel *goquery.Selection) {
				tagName := goquery.NodeName(sel)
				var elemType entities.ElementType
				switch tagName {
				case "h1":
					elemType = entities.ElementTypeH1
				case "h2":
					elemType = entities.ElementTypeH2
				case "h3":
					elemType = entities.ElementTypeH3
				case "p":
					elemType = entities.ElementTypeParagraph
				}

				// Создаем RichText элемент
				text := strings.TrimSpace(sel.Text())
				richText := entities.RichText{
					PlainText: text,
					Meta:      []entities.Meta{},
				}

				// Проверяем форматирование
				sel.Children().Each(func(i int, child *goquery.Selection) {
					childName := goquery.NodeName(child)
					childText := child.Text()
					offset := strings.Index(text, childText)
					length := len(childText)

					meta := entities.Meta{
						Offset: offset,
						Length: length,
					}

					switch childName {
					case "b", "strong":
						meta.Weight = "bold"
					case "i", "em":
						meta.Italic = true
					case "u":
						meta.Underline = true
					case "strike", "s":
						meta.Strikethrough = true
					case "a":
						if href, exists := child.Attr("href"); exists {
							meta.Link = href
						}
					}

					if offset >= 0 && length > 0 {
						richText.Meta = append(richText.Meta, meta)
					}
				})

				contentElements = append(contentElements, entities.TextElement{
					Type:     elemType,
					RichText: richText,
				})
			})

			// Парсим изображения
			s.Find("img").Each(func(i int, sel *goquery.Selection) {
				src, _ := sel.Attr("src")
				widthStr, _ := sel.Attr("width")
				align, _ := sel.Attr("align")

				width := 0
				if widthStr != "" {
					fmt.Sscanf(widthStr, "%d", &width)
				}

				if align == "" {
					align = "left"
				}

				contentElements = append(contentElements, entities.ImgElement{
					Type:  entities.ElementTypeImage,
					Src:   src,
					Width: width,
					Align: align,
				})
			})
		}
	})

	// Создаем заметку
	note := &entities.Note{
		Source:         entities.Source(url),
		Title:          doc.Find("title").Text(),
		ContentDecoded: contentElements,
	}

	// Сериализуем content в JSON
	// contentJSON, err := json.Marshal(contentElements)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal content: %v", err)
	// }
	// note.Content = json.RawMessage(contentJSON)

	return note, nil
}
