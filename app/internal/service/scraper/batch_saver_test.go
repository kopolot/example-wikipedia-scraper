package scraper

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/testutils"
	mockRepo "example-wikipedia-scraper/internal/testutils/repository"
	"testing"

	"github.com/stretchr/testify/mock"
)

var (
	loggerSaver       *testutils.MockLogger
	pageFactory       *testutils.MockPageFactory
	repository        *mockRepo.MockPageRepository
	failedPages       chan *dto.UnprocessedPageDTO
	queryBuilderSaver *testutils.MockQueryBuilder

	nilSlice = []interface{}(nil)
)

func setUpBatchSaver(t *testing.T) *BatchSaver {
	t.Helper()
	loggerSaver = &testutils.MockLogger{}
	pageFactory = &testutils.MockPageFactory{}
	repository = &mockRepo.MockPageRepository{}
	failedPages = make(chan *dto.UnprocessedPageDTO, 10)
	queryBuilderSaver = &testutils.MockQueryBuilder{}
	batchSaver := NewBatchSaver(loggerSaver, pageFactory, repository, failedPages, queryBuilderSaver)
	return batchSaver
}

func TestSaveBatch_NoSamePages(t *testing.T) {
	batchSaver := setUpBatchSaver(t)
	batch := []*dto.PageDTO{
		{
			Title: "Page 1",
		},
		{
			Title: "Page 2",
		},
	}

	pages := []*model.Page{
		{Title: "Page 1", HashKey: "1", Notified: false},
		{Title: "Page 2", HashKey: "2", Notified: false},
	}

	pageFactory.On("CreateFromDTO", batch[0]).Return(pages[0], nil)
	pageFactory.On("CreateFromDTO", batch[1]).Return(pages[1], nil)

	result := make([]string, 0)
	queryBuilderSaver.
		On("Table", "pages", nilSlice).
		On("Select", "hash_key", nilSlice).
		On("Where", "notified = true AND hash_key IN ?", []interface{}{[]string{"1", "2"}}).
		On("Find", &result, nilSlice).
		Return(nil)

	repository.On("Upsert", pages[0], []string{"url"}).Return(nil)
	repository.On("Upsert", pages[1], []string{"url"}).Return(nil)

	batchSaver.SaveBatch(batch)

	pageFactory.AssertExpectations(t)
	queryBuilderSaver.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func TestSaveBatch_HasSamePage(t *testing.T) {
	batchSaver := setUpBatchSaver(t)
	batch := []*dto.PageDTO{
		{
			Title: "Page 1",
		},
		{
			Title: "Page 2",
		},
	}

	pages := []*model.Page{
		{Title: "Page 1", HashKey: "1", Notified: false},
		{Title: "Page 2", HashKey: "2", Notified: false},
	}

	pageFactory.On("CreateFromDTO", batch[0]).Return(pages[0], nil)
	pageFactory.On("CreateFromDTO", batch[1]).Return(pages[1], nil)

	result := make([]string, 0)
	queryBuilderSaver.
		On("Table", "pages", nilSlice).
		On("Select", "hash_key", nilSlice).
		On("Where", "notified = true AND hash_key IN ?", []interface{}{[]string{"1", "2"}}).
		On("Find", &result, nilSlice).
		Run(func(args mock.Arguments) {
			res := args.Get(0).(*[]string)
			*res = []string{"1"}
		}).
		Return(nil)

	pages[0].Notified = true

	repository.On("Upsert", pages[0], []string{"url"}).Return(nil)
	repository.On("Upsert", pages[1], []string{"url"}).Return(nil)

	batchSaver.SaveBatch(batch)

	pageFactory.AssertExpectations(t)
	queryBuilderSaver.AssertExpectations(t)
	repository.AssertExpectations(t)
}
