package utils

import (
	"example-wikipedia-scraper/internal/factory"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func InsertTestPage(t *testing.T, page *model.Page) *model.Page {
	err := repository.NewPageRepository().Create(page)
	require.NoError(t, err)
	return page
}

func DeleteTestPage(t *testing.T, pageID uint) {
	err := repository.NewPageRepository().Delete(pageID)
	require.NoError(t, err)
}

func GenerateTestPage(t *testing.T, count int) []*model.Page {
	pages := make([]*model.Page, 0, count)
	for i := range count {
		j := i + 1
		page := &model.Page{
			SiteName:   "TestSite",
			URL:        "https://example.com/page/" + strconv.Itoa(j),
			Content:    "This is the content of test page " + strconv.Itoa(j),
			Title:      "Test Page " + strconv.Itoa(j),
			TextField1: "Additional info 1 for page " + strconv.Itoa(j),
			TextField2: "Additional info 2 for page " + strconv.Itoa(j),
			TextField3: "Additional info 3 for page " + strconv.Itoa(j),
			ExternalID: uuid.NewString(),
			Notified:   false,
		}
		page.HashKey = factory.GenerateHashKey(page)
		pages = append(pages, page)
	}
	return pages
}

func CreateTestPages(t *testing.T, count int) []*model.Page {
	pages := GenerateTestPage(t, count)
	err := repository.NewPageRepository().CreateInBatches(pages, uint(count))
	require.NoError(t, err)
	return pages
}
