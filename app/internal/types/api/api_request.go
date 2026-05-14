package api

import "example-wikipedia-scraper/internal/model"

type ApiRequest struct {
	Method      string              `json:"method"`
	URL         string              `json:"url"`
	Body        string              `json:"body"`
	ClientIP    string              `json:"client_ip"`
	Headers     map[string]string   `json:"headers"`
	QueryParams map[string][]string `json:"queryParams"`
	PathParams  map[string]string   `json:"pathParams"`
	User        *model.User
}
