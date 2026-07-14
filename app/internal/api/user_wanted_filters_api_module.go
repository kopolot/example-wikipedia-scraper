package api

import (
	"encoding/json"
	"example-wikipedia-scraper/internal/interfaces/api"
	"example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/interfaces/service"
	"example-wikipedia-scraper/internal/model"
	httpPolicy "example-wikipedia-scraper/internal/policy/http"
	types "example-wikipedia-scraper/internal/types/api"
	"strconv"
)

type UserWantedFiltersApiModule struct {
	api                       api.ApiInterface
	userWantedPagesFilterRepo repository.UserWantedPagesFilterRepositoryInterface
	pageFilterService         service.PageFilterServiceInterface
	subscriptionService       service.SubscriptionServiceInterface
}

func NewUserWantedFiltersApiModule(
	api api.ApiInterface,
	userWantedPagesFilterRepo repository.UserWantedPagesFilterRepositoryInterface,
	pageFilterService service.PageFilterServiceInterface,
	subscriptionService service.SubscriptionServiceInterface,
) *UserWantedFiltersApiModule {
	return &UserWantedFiltersApiModule{
		api:                       api,
		userWantedPagesFilterRepo: userWantedPagesFilterRepo,
		pageFilterService:         pageFilterService,
		subscriptionService:       subscriptionService,
	}
}

func (a *UserWantedFiltersApiModule) GetRoutes() []*types.Route {
	policy := httpPolicy.NewUserWantedPagesFiltersPolicy(a.subscriptionService)
	return []*types.Route{
		{
			Method:      "GET",
			Path:        "/",
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware},
			Handler:     a.getUserWantedPagesFilters,
		},
		{
			Method:      "GET",
			Path:        "/filtered_pages",
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, makeMiddlewareFromPolicy(policy.CanGetFilteredPages, a.userWantedPagesFilterRepo)},
			Handler:     a.getFilteredPages,
		},
		{
			Method:      "GET",
			Path:        "/pages_by_filters",
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, makeMiddlewareFromPolicy(policy.CanGetFilteredPages, a.userWantedPagesFilterRepo)},
			Handler:     a.getPagesByFilters,
		},
		{
			Method:      "POST",
			Path:        "/",
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, makeMiddlewareFromPolicy(policy.CanCreate, a.userWantedPagesFilterRepo)},
			Handler:     a.createUserWantedPagesFilter,
		},
		{
			Method:      "DELETE",
			Path:        "/:id",
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, makeMiddlewareFromPolicy(policy.CanDelete, a.userWantedPagesFilterRepo)},
			Handler:     a.deleteUserWantedPagesFilter,
		},
		{
			Method:      "PUT",
			Path:        "/:id",
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, makeMiddlewareFromPolicy(policy.CanUpdate, a.userWantedPagesFilterRepo)},
			Handler:     a.updateUserWantedPagesFilter,
		},
		{
			Method:      "GET",
			Path:        "/:id",
			Middlewares: []types.Middleware{a.api.AuthenticateMiddleware, makeMiddlewareFromPolicy(policy.CanView, a.userWantedPagesFilterRepo)},
			Handler:     a.getUserWantedPagesFilter,
		},
	}
}

func (a *UserWantedFiltersApiModule) GetRoutePrefix() string {
	return "user_wanted_filters"
}

func (a *UserWantedFiltersApiModule) getUserWantedPagesFilters(request *types.ApiRequest) *types.ApiResponse {
	user := request.User
	var userFilters []*model.UserWantedPagesFilter
	err := a.userWantedPagesFilterRepo.GetQueryBuilder().Unscoped().Where("user_id = ?", user.ID).Find(&userFilters)
	if err != nil {
		return InternalErrorResponse()
	}
	return OkResponse(userFilters)
}

func (a *UserWantedFiltersApiModule) createUserWantedPagesFilter(request *types.ApiRequest) *types.ApiResponse {
	user := request.User
	var criteria model.UserWantedPageCriteria
	if err := json.Unmarshal([]byte(request.Body), &criteria); err != nil {
		return BadRequestResponse()
	}
	userFilter := &model.UserWantedPagesFilter{
		FilterData: criteria,
		UserID:     user.ID,
	}
	if err := a.userWantedPagesFilterRepo.Create(userFilter); err != nil {
		return InternalErrorResponse()
	}
	return CreatedResponse(userFilter)
}

func (a *UserWantedFiltersApiModule) deleteUserWantedPagesFilter(request *types.ApiRequest) *types.ApiResponse {
	stringId := request.PathParams["id"]
	id, err := strconv.ParseUint(stringId, 10, 32)
	if err != nil {
		return BadRequestResponseWithMsg("Invalid filter ID")
	}
	err = a.userWantedPagesFilterRepo.GetQueryBuilder().
		Unscoped().
		Where("id = ? AND user_id = ?", uint(id), request.User.ID).
		Delete(&model.UserWantedPagesFilter{})
	if err != nil {
		return InternalErrorResponse()
	}
	return OkResponse(nil)
}

func (a *UserWantedFiltersApiModule) updateUserWantedPagesFilter(request *types.ApiRequest) *types.ApiResponse {
	stringId := request.PathParams["id"]
	id, err := strconv.ParseUint(stringId, 10, 32)
	if err != nil {
		return BadRequestResponseWithMsg("Invalid filter ID")
	}
	var criteria model.UserWantedPageCriteria
	if err := json.Unmarshal([]byte(request.Body), &criteria); err != nil {
		return BadRequestResponse()
	}
	userFilter, err := a.userWantedPagesFilterRepo.GetByID(uint(id))
	if err != nil {
		return NotFoundResponse("Filter not found")
	}
	userFilter.FilterData = criteria
	if err := a.userWantedPagesFilterRepo.Update(userFilter); err != nil {
		return InternalErrorResponse()
	}
	return OkResponse(userFilter)
}

func (a *UserWantedFiltersApiModule) getUserWantedPagesFilter(request *types.ApiRequest) *types.ApiResponse {
	stringId := request.PathParams["id"]
	id, err := strconv.ParseUint(stringId, 10, 32)
	if err != nil {
		return BadRequestResponse()
	}
	var userFilter *model.UserWantedPagesFilter
	err = a.userWantedPagesFilterRepo.GetQueryBuilder().Unscoped().First(&userFilter, uint(id))
	if err != nil {
		return NotFoundResponse("User wanted pages filter not found")
	}
	return OkResponse(userFilter)
}

func (a *UserWantedFiltersApiModule) getFilteredPages(request *types.ApiRequest) *types.ApiResponse {
	_, maxRecords, page := a.getFilteredPagesRequestData(request)
	filters, err := a.userWantedPagesFilterRepo.FindBy("user_id = ?", request.User.ID)
	if err != nil {
		return InternalErrorResponse()
	}
	if len(filters) == 0 {
		return OkResponse([]*model.Page{})
	}
	filterIDs := make([]uint, 0, len(filters))
	for _, filter := range filters {
		filterIDs = append(filterIDs, filter.ID)
	}
	pages, err := a.pageFilterService.GetPagesByFilters(filterIDs, maxRecords, page)
	if err != nil {
		return InternalErrorResponse()
	}
	return OkResponse(pages)
}

func (a *UserWantedFiltersApiModule) getFilteredPagesRequestData(request *types.ApiRequest) (uint, uint, uint) {
	maxRecords := uint(100)
	page := uint(1)
	if len(request.QueryParams["page"]) > 0 {
		if p, err := strconv.Atoi(request.QueryParams["page"][0]); err == nil && p > 0 {
			page = uint(p)
		}
	}
	if len(request.QueryParams["max_records"]) > 0 {
		if mr, err := strconv.Atoi(request.QueryParams["max_records"][0]); err == nil && mr > 0 && mr <= 1000 {
			maxRecords = uint(mr)
		}
	}
	return request.User.ID, maxRecords, page
}

func (a *UserWantedFiltersApiModule) getPagesByFilters(request *types.ApiRequest) *types.ApiResponse {
	userId := request.User.ID
	stringIds := request.QueryParams["id"]
	ids := make([]uint, 0, len(stringIds))
	for _, idStr := range stringIds {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			return BadRequestResponse()
		}
		ids = append(ids, uint(id))
	}
	filters, err := a.userWantedPagesFilterRepo.FindBy("id IN ? AND user_id = ?", ids, userId)
	if err != nil {
		return InternalErrorResponse()
	}
	if len(filters) == 0 {
		return OkResponse([]*model.Page{})
	}
	foundIDs := make([]uint, 0, len(filters))
	for _, filter := range filters {
		foundIDs = append(foundIDs, filter.ID)
	}
	pages, err := a.pageFilterService.GetPagesByFilters(foundIDs, 100, 1)
	if err != nil {
		return InternalErrorResponse()
	}
	return OkResponse(pages)
}
