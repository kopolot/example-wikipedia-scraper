package api

type Route struct {
	Method      string
	Path        string
	Action      func()
	Handler     func(request *ApiRequest) *ApiResponse
	Middlewares []Middleware
}

type Middleware func(request *ApiRequest) *ApiResponse
