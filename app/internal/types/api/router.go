package api

type Route struct {
	Method           string
	Path             string
	Handler          func(request *ApiRequest) *ApiResponse
	Middlewares      []Middleware
	AfterMiddlewares []AfterMiddleware
}

type Middleware func(request *ApiRequest) *ApiResponse

type AfterMiddleware func(request *ApiRequest, response *ApiResponse)
