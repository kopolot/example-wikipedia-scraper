package factory

import (
	"crypto/sha256"
	"errors"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
	"fmt"
	"html"
	"regexp"
	"strings"
)

var (
	errEmptyMatchValue    = fmt.Errorf("value must be non-empty")
	errEmptyMatchProperty = fmt.Errorf("property must be non-empty")
	stripHTMLRegex        = regexp.MustCompile(`<[^>]*>`)
)

type PageFactory struct {
}

func NewPageFactory() *PageFactory {
	return &PageFactory{}
}

func (f *PageFactory) CreateFromDTO(pageDTO *dto.PageDTO) (*model.Page, error) {
	if pageDTO == nil {
		return nil, errors.New("pageDTO cannot be nil")
	}
	page := &model.Page{
		Title:      pageDTO.Title,
		URL:        pageDTO.URL,
		SiteName:   pageDTO.SiteName,
		Content:    pageDTO.Content,
		TextField1: pageDTO.TextField1,
		TextField2: pageDTO.TextField2,
		TextField3: pageDTO.TextField3,
		ExternalID: pageDTO.ExternalID,
	}
	page.HashKey = GenerateHashKey(page)
	return page, nil
}

func GenerateHashKey(page *model.Page) string {

	components := []string{
		page.TextField1,
		page.TextField2,
		page.TextField3,
	}
	stringKey := strings.ToLower(strings.Join(components, "|"))
	hashKey := sha256.Sum256([]byte(stringKey))
	return fmt.Sprintf("%x", hashKey)
}

func stripHTML(input string) string {
	return stripHTMLRegex.ReplaceAllString(input, "")
}

func normalizeText(input string) string {
	clean := html.UnescapeString(input)
	clean = stripHTML(clean)
	clean = strings.ToLower(clean)
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "\n", "")
	clean = strings.ReplaceAll(clean, "\t", "")
	return clean
}
