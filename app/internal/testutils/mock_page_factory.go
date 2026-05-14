package testutils

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockPageFactory struct {
	mock.Mock
}

func (m *MockPageFactory) CreateFromDTO(dto *dto.PageDTO) (*model.Page, error) {
	args := m.Called(dto)
	return args.Get(0).(*model.Page), args.Error(1)
}
