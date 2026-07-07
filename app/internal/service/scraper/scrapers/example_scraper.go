package scrapers

import (
	"context"
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
	types "example-wikipedia-scraper/internal/types/scraper"
	"fmt"
	"log"
	"strconv"
	"time"
)

type ExampleScraper struct {
	url           string
	scrapeOptions *types.ScrapeOptions
}

func NewExampleScraper(url string) *ExampleScraper {
	return &ExampleScraper{
		url: url,
	}
}

func (s *ExampleScraper) GetScrapeOptions() types.ScrapeOptions {
	if s.scrapeOptions == nil {
		return types.ScrapeOptions{}
	}
	return *s.scrapeOptions
}

func (s *ExampleScraper) InitScraper(opts ...types.ScrapeOption) error {
	s.scrapeOptions = types.ApplyOptions(opts...)
	return nil
}

func (s *ExampleScraper) GetName() string {
	return "example"
}

func (s *ExampleScraper) GetURL() string {
	return s.url
}

func (s *ExampleScraper) ScrapeAsync(channels *types.ScrapeChannels) error {
	pageQueue := channels.PageQueue
	options := s.GetScrapeOptions()

	if options.MaxItems <= 0 {
		options.MaxItems = 10
	}
	timeout := 30 * time.Second
	if options.Timeout != 0 {
		timeout = options.Timeout
	}

	log.Printf("Example scraper: starting async scraping (maxItems: %d, timeout: %v)", options.MaxItems, timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for i := 0; i < options.MaxItems; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("scraping timeout after %v", timeout)
		default:

			time.Sleep(10 * time.Millisecond)

			page := s.createExamplePageDTO(i)

			select {
			case pageQueue <- page:
				log.Printf("Example scraper: sent page %d to queue", i)
			case <-ctx.Done():
				return fmt.Errorf("timeout while sending page %d to queue", i)
			}
		}
	}

	log.Printf("Example scraper: completed async scraping of %d pages", options.MaxItems)
	return nil
}

func (s *ExampleScraper) ScrapeSync(opts ...types.ScrapeOption) ([]model.Page, error) {
	options := types.ApplyOptions(opts...)

	if options.MaxItems <= 0 {
		options.MaxItems = 10
	}
	timeout := 30 * time.Second
	if options.Timeout != 0 {
		timeout = options.Timeout
	}

	log.Printf("Example scraper: starting sync scraping (maxItems: %d, timeout: %v)", options.MaxItems, timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pages := make([]model.Page, 0, options.MaxItems)

	for i := 0; i < options.MaxItems; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("scraping timeout after %v (scraped %d/%d pages)", timeout, len(pages), options.MaxItems)
		default:

			time.Sleep(10 * time.Millisecond)

			page := s.createExamplePage(i)
			pages = append(pages, page)
			log.Printf("Example scraper: scraped page %d", i)
		}
	}

	log.Printf("Example scraper: completed sync scraping of %d pages", len(pages))
	return pages, nil
}

func (s *ExampleScraper) ScrapePageData(pageURL string) (*dto.PageDTO, error) {
	log.Printf("Example scraper: scraping page data from %s", pageURL)

	time.Sleep(50 * time.Millisecond)

	pageID := 0
	if len(pageURL) > 0 {

		for i := len(pageURL) - 1; i >= 0; i-- {
			if pageURL[i] >= '0' && pageURL[i] <= '9' {
				if id, err := strconv.Atoi(string(pageURL[i])); err == nil {
					pageID = id
					break
				}
			}
		}
	}

	page := s.createDetailedExamplePageDTO(pageID, pageURL)
	log.Printf("Example scraper: scraped detailed page data for page %d", pageID)

	return page, nil
}

func (s *ExampleScraper) ValidatePage(page *model.Page) (*dto.PageDTO, error) {
	if page == nil {
		return nil, nil
	}

	log.Printf("Example scraper: validating page: %s", page.Title)

	if page.Title == "" {
		return nil, fmt.Errorf("invalid page: title is empty")
	}
	if page.URL == "" {
		return nil, fmt.Errorf("invalid page: URL is empty")
	}
	if page.SiteName != "example" {
		return nil, fmt.Errorf("invalid page: wrong site name '%s', expected 'example'", page.SiteName)
	}
	validated := &dto.PageDTO{
		Title:      page.Title,
		URL:        page.URL,
		SiteName:   page.SiteName,
		Content:    page.Content,
		TextField1: page.TextField1,
		TextField2: page.TextField2,
		TextField3: page.TextField3,
	}

	log.Printf("Example scraper: page validated successfully")
	return validated, nil
}

func (s *ExampleScraper) createExamplePageDTO(index int) *dto.PageDTO {
	return &dto.PageDTO{
		Title:      fmt.Sprintf("Example Page %d", index),
		URL:        fmt.Sprintf("%s/page/%d", s.url, index),
		SiteName:   "example",
		Content:    fmt.Sprintf("This is the content of example page %d.", index),
		TextField1: fmt.Sprintf("TextField1 for page %d", index),
		TextField2: fmt.Sprintf("TextField2 for page %d", index),
		TextField3: fmt.Sprintf("TextField3 for page %d", index),
	}
}

func (s *ExampleScraper) createExamplePage(index int) model.Page {
	return model.Page{
		Title:      fmt.Sprintf("Example Page %d", index),
		URL:        fmt.Sprintf("%s/page/%d", s.url, index),
		SiteName:   "example",
		Content:    fmt.Sprintf("This is the content of example page %d.", index),
		TextField1: fmt.Sprintf("TextField1 for page %d", index),
		TextField2: fmt.Sprintf("TextField2 for page %d", index),
		TextField3: fmt.Sprintf("TextField3 for page %d", index),
		// Additional fields to simulate a real page
	}
}

func (s *ExampleScraper) createDetailedExamplePageDTO(pageID int, pageURL string) *dto.PageDTO {
	description := fmt.Sprintf("This is an example page description for page %d. Very detailed and informative.", pageID)
	return &dto.PageDTO{
		Title:      fmt.Sprintf("Example Page %d", pageID),
		URL:        pageURL,
		SiteName:   "example",
		Content:    description,
		TextField1: fmt.Sprintf("TextField1 for page %d", pageID),
		TextField2: fmt.Sprintf("TextField2 for page %d", pageID),
		TextField3: fmt.Sprintf("TextField3 for page %d", pageID),
		// Additional fields could be added here as needed
	}
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
