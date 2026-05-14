package testutils

import (
	"example-wikipedia-scraper/internal/types/api"

	"github.com/stretchr/testify/mock"
)

type MockApiModule struct {
	mock.Mock
}

func (m *MockApiModule) GetRoutes() []*api.Route {
	args := m.Called()
	return args.Get(0).([]*api.Route)
}

func (m *MockApiModule) GetRoutePrefix() string {
	args := m.Called()
	return args.String(0)
}
