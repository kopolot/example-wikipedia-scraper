package api

import (
	"example-wikipedia-scraper/internal/interfaces/api"
	"example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/model"
	types "example-wikipedia-scraper/internal/types/api"
	"math"
	"net/http"
	"strconv"
)

type PageApiModule struct {
	api      api.ApiInterface
	pageRepo repository.PageRepositoryInterface
}

func NewPageApiModule(api api.ApiInterface, pageRepo repository.PageRepositoryInterface) *PageApiModule {
	return &PageApiModule{
		api:      api,
		pageRepo: pageRepo,
	}
}

func (m *PageApiModule) GetRoutes() []*types.Route {
	return []*types.Route{
		{
			Method:  http.MethodGet,
			Path:    "/pages",
			Handler: m.getPages,
		},
	}
}

func (m *PageApiModule) GetRoutePrefix() string {
	return "page"
}

var allowedSortFields = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
	"title":      "title",
}

func (m *PageApiModule) getPages(request *types.ApiRequest) *types.ApiResponse {
	page := uint(1)
	limit := uint(100)
	sortField := "created_at"
	sortOrder := "desc"

	if val, exists := request.QueryParams["page"]; exists && len(val) > 0 {
		if pageParsed, err := strconv.ParseUint(val[0], 10, 32); err == nil && pageParsed > 0 {
			page = uint(pageParsed)
		}
	}
	if val, exists := request.QueryParams["limit"]; exists && len(val) > 0 {
		if limitParsed, err := strconv.ParseUint(val[0], 10, 32); err == nil && limitParsed > 0 {
			limit = uint(limitParsed)
		}
	}
	if val, exists := request.QueryParams["sort"]; exists && len(val) > 0 {
		if col, ok := allowedSortFields[val[0]]; ok {
			sortField = col
		}
	}
	if val, exists := request.QueryParams["order"]; exists && len(val) > 0 && (val[0] == "asc" || val[0] == "desc") {
		sortOrder = val[0]
	}

	qb := m.pageRepo.GetQueryBuilder().Order(sortField + " " + sortOrder)

	// if val, exists := request.QueryParams["price_max"]; exists && len(val) > 0 {
	// 	if v, err := strconv.ParseUint(val[0], 10, 32); err == nil {
	// 		qb = qb.Where("price <= ?", v)
	// 	}
	// }

	offset := (page - 1) * limit
	var pageRecords []*model.Page
	if err := qb.Limit(int(limit)).Offset(int(offset)).Find(&pageRecords); err != nil {
		return InternalErrorResponse()
	}

	var totalCount int64
	countQb := m.pageRepo.GetQueryBuilder().Model(&model.Page{})

	// if val, exists := request.QueryParams["price_max"]; exists && len(val) > 0 {
	// 	if v, err := strconv.ParseUint(val[0], 10, 32); err == nil {
	// 		countQb = countQb.Where("price <= ?", v)
	// 	}
	// }
	if err := countQb.Count(&totalCount); err != nil {
		return InternalErrorResponse()
	}

	maxPage := uint(math.Ceil(float64(totalCount) / float64(limit)))
	if maxPage == 0 {
		maxPage = 1
	}
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
			Data:    map[string]any{"pageRecords": pageRecords, "maxPage": maxPage},
		},
	}
}
