package service

import (
	repoInterface "example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/types/service"
)

type PageFilterService struct {
	userWantedPagesFiltersRepo repoInterface.UserWantedPagesFilterRepositoryInterface
	pageRepo                   repoInterface.PageRepositoryInterface
}

func NewPageFilterService(
	userWantedPagesFiltersRepo repoInterface.UserWantedPagesFilterRepositoryInterface,
	pageRepo repoInterface.PageRepositoryInterface,
) *PageFilterService {
	return &PageFilterService{
		userWantedPagesFiltersRepo: userWantedPagesFiltersRepo,
		pageRepo:                   pageRepo,
	}
}

func (s *PageFilterService) GetPagesByFilters(ids []uint, limit, page uint) ([]*model.Page, error) {
	filters, err := s.userWantedPagesFiltersRepo.FindBy("id IN ?", ids)
	if err != nil {
		return nil, err
	}
	var filterQuery service.FilterQuery
	for _, filter := range filters {
		currentFilterQuery := s.BuildFilterQuery(filter.FilterData)
		filterQuery.CombineQuerys(currentFilterQuery)
	}
	query, args := filterQuery.GetQuery()
	if query == "" {
		return []*model.Page{}, nil
	}
	args = append([]any{query}, args...)
	pages, err := s.pageRepo.FindByWithLimitPage(args, limit, page)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (s *PageFilterService) BuildFilterQuery(filter model.UserWantedPageCriteria) *service.FilterQuery {
	filterQuery := &service.FilterQuery{}
	if len(filter.SiteNames) > 0 {
		filterQuery.AddCondition("site_name IN ?", []string(filter.SiteNames))
	}
	for _, keyword := range filter.Keywords {
		like := "%" + keyword + "%"
		filterQuery.AddCondition("(title ILIKE ? OR content ILIKE ?)", like, like)
	}
	if filter.TitleContains != "" {
		filterQuery.AddCondition("title ILIKE ?", "%"+filter.TitleContains+"%")
	}
	return filterQuery
}

func (s *PageFilterService) FindMatchingFiltersForPage(page model.Page) ([]*model.UserWantedPagesFilter, error) {
	reverseFilterQuery := s.BuildReverseFilterQuery(page)
	query, args := reverseFilterQuery.GetQuery()
	if query == "" {
		return []*model.UserWantedPagesFilter{}, nil
	}
	args = append([]any{query}, args...)
	filters, err := s.userWantedPagesFiltersRepo.FindBy(args...)
	if err != nil {
		return nil, err
	}
	return filters, nil
}

func (s *PageFilterService) BuildReverseFilterQuery(page model.Page) *service.ReverseFilterQuery {
	reverseFilterQuery := &service.ReverseFilterQuery{}
	reverseFilterQuery.AddCondition(
		"(site_names IS NULL OR cardinality(site_names) = 0 OR ? = ANY(site_names))",
		page.SiteName,
	)
	reverseFilterQuery.AddCondition(
		"(keywords IS NULL OR cardinality(keywords) = 0 OR NOT EXISTS (SELECT 1 FROM unnest(keywords) AS kw WHERE (? NOT ILIKE '%' || kw || '%' AND ? NOT ILIKE '%' || kw || '%')))",
		page.Title,
		page.Content,
	)
	reverseFilterQuery.AddCondition(
		"(title_contains IS NULL OR title_contains = '' OR ? ILIKE '%' || title_contains || '%')",
		page.Title,
	)
	return reverseFilterQuery
}
