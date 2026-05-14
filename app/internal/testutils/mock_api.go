package testutils

import (
	"example-wikipedia-scraper/internal/interfaces"
	apiiface "example-wikipedia-scraper/internal/interfaces/api"
	"example-wikipedia-scraper/internal/interfaces/auth"
	"example-wikipedia-scraper/internal/types/api"

	"github.com/stretchr/testify/mock"
)

type MockApi struct {
	mock.Mock
	Logger      interfaces.LoggerInterface
	AuthManager auth.AuthManagerInterface
	AuthResp    *api.ApiResponse
}

func (m *MockApi) AuthenticateMiddleware(request *api.ApiRequest) *api.ApiResponse {
	return m.AuthResp
}
func (m *MockApi) GetLogger() interfaces.LoggerInterface {
	return m.Logger
}
func (m *MockApi) LoadModule(module apiiface.ApiModuleInterface) {}
func (m *MockApi) GetAuthManager() auth.AuthManagerInterface {
	return m.AuthManager
}
