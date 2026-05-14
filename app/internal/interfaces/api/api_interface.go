package api

import (
	"example-wikipedia-scraper/internal/interfaces"
	authType "example-wikipedia-scraper/internal/interfaces/auth"
	types "example-wikipedia-scraper/internal/types/api"
)

const (
	MsgAuthHeaderMissing = "Authorization header missing"
)

type ApiInterface interface {
	AuthenticateMiddleware(request *types.ApiRequest) *types.ApiResponse
	GetLogger() interfaces.LoggerInterface
	LoadModule(module ApiModuleInterface)
	GetAuthManager() authType.AuthManagerInterface
}
