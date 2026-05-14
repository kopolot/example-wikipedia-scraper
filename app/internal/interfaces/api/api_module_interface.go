package api

import (
	types "example-wikipedia-scraper/internal/types/api"
)

type ApiModuleInterface interface {
	GetRoutes() []*types.Route
	GetRoutePrefix() string
}
