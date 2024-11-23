package main

import (
	"regexp"

	"github.com/Corray333/keep_it/parsers/vk/internal/app"
	"github.com/Corray333/keep_it/parsers/vk/internal/entities"
)

func main() {
	app.New().Run()
	// fmt.Println(parseRichText("Test [https://vk.com/studiyagrafit?w=wall-57775506_20297|сессии говна] на на [https://vk.com/studiyagrafit?w=wall-57775506_20297|сессии говна] на"))
}

func parseRichText(input string) (result entities.RichText) {
	var meta []entities.Meta
	offsetCorrection := 0

	// Регулярное выражение для ссылок и текста
	linkRegex := regexp.MustCompile(`\[(https?://[^\|]+)\|([^\]]+)\]`)
	linkMatches := linkRegex.FindAllStringSubmatchIndex(input, -1)
	for i := range linkMatches {
		if i == 0 && linkMatches[i][0] > 0 {
			result.PlainText += input[offsetCorrection:linkMatches[i][0]]
		}
		link := input[linkMatches[i][2]:linkMatches[i][3]]
		text := input[linkMatches[i][4]:linkMatches[i][5]]
		meta = append(meta, entities.Meta{
			Offset: len([]rune(result.PlainText)),
			Length: len(text),
			Link:   link,
		})

		result.PlainText += text
		if i < len(linkMatches)-1 {
			result.PlainText += input[linkMatches[i][1]:linkMatches[i+1][0]]
		} else {
			result.PlainText += input[linkMatches[i][1]:]
		}
	}

	// После обработки метаинформации оставшийся текст является plainText
	result.Meta = meta
	return result
}

// Test сессии говна на на сессии говна на
