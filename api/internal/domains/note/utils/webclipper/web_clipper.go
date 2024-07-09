package webclipper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

func GetMD(url string) string {

	html, err := fetchHTML(url)
	if err != nil {
		log.Fatal("Error fetching HTML:", err)
	}

	markdownArticles, err := extractAndConvertArticles(html)
	if err != nil {
		log.Fatal("Error extracting and converting articles:", err)
	}

	result := ""

	for _, markdown := range markdownArticles {
		result += markdown + "\n"
	}
	return result
}

func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch URL: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	html, err := doc.Html()
	if err != nil {
		return "", err
	}

	return html, nil
}

func extractAndConvertArticles(html string) ([]string, error) {
	converter := md.NewConverter("", true, nil)
	var markdownArticles []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	doc.Find("article").Find("p,h1,h2,h3,img").Each(func(i int, s *goquery.Selection) {
		articleHtml, err := s.Html()
		if err != nil {
			log.Println("Error extracting article HTML:", err)
			return
		}

		markdown, err := converter.ConvertString(articleHtml)
		if err != nil {
			log.Println("Error converting article to markdown:", err)
			return
		}

		markdownArticles = append(markdownArticles, markdown)
	})

	return markdownArticles, nil
}
