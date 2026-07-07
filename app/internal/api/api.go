package api

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/interfaces"
	apiInterfaces "example-wikipedia-scraper/internal/interfaces/api"
	authType "example-wikipedia-scraper/internal/interfaces/auth"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/service"
	"example-wikipedia-scraper/internal/service/mailer"
	types "example-wikipedia-scraper/internal/types/api"
	pkgRepo "example-wikipedia-scraper/pkg/repository"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	contextRequestKey = "apiRequest"
)

type Api struct {
	logger      interfaces.LoggerInterface
	authManager authType.AuthManagerInterface
	config      config.ConfigInterface
	engine      *gin.Engine
	routes      map[string]*types.Route
}

func NewApi(config config.ConfigInterface, logger interfaces.LoggerInterface, authManager authType.AuthManagerInterface) *Api {
	if config.GetApiConfig().Debug {
		gin.SetMode(gin.DebugMode)
		logger.Info("API running in debug mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.GetLogWriter())
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, logger.GetLogWriter())
	engine := gin.Default()
	api := &Api{
		config:      config,
		logger:      logger,
		engine:      engine,
		routes:      make(map[string]*types.Route),
		authManager: authManager,
	}
	healthCheckRoute := &types.Route{
		Method:  http.MethodGet,
		Path:    "/health",
		Handler: api.healthCheck,
	}
	api.RegisterRoutes([]*types.Route{healthCheckRoute}, "")
	return api
}

func (a *Api) LoadModules() {
	a.loadUserApiModule()
	a.loadPageApiModule()
	a.loadUserWantedFiltersApiModule()
}

func (a *Api) loadUserApiModule() {
	userRepo := repository.NewUserRepository()
	mailerService := mailer.NewMailer(a.config.GetMailerConfig(), a.logger)
	userService := service.NewUserService(userRepo, mailerService, a.config)
	userApiModule := NewUserApiModule(userRepo, userService, a)
	a.LoadModule(userApiModule)
}

func (a *Api) loadPageApiModule() {
	pageRepo := repository.NewPageRepository()
	pageApiModule := NewPageApiModule(a, pageRepo)
	a.LoadModule(pageApiModule)
}

func (a *Api) loadUserWantedFiltersApiModule() {
	userWantedPagesFilterRepo := repository.NewUserWantedPagesFilterRepository()
	pageRepo := repository.NewPageRepository()
	pageFilterService := service.NewPageFilterService(userWantedPagesFilterRepo, pageRepo)
	userWantedFiltersApiModule := NewUserWantedFiltersApiModule(a, userWantedPagesFilterRepo, pageFilterService)
	a.LoadModule(userWantedFiltersApiModule)
}

func (a *Api) RegisterRoutes(routes []*types.Route, prefix string) {
	for _, route := range routes {
		if prefix != "" {
			route.Path = strings.TrimPrefix(route.Path, "")
			route.Path = strings.TrimPrefix(prefix, "") + "" + route.Path
		}
		route.Method = strings.ToUpper(route.Method)
		a.routes[route.Path+"_"+route.Method] = route
	}
}

func (a *Api) SetupRoutes() {
	a.logger.Info("Setting up routes")
	api := a.engine
	// api := a.engine.Group("/api")
	api.Use(a.recoverFromPanicMiddleware)
	for _, route := range a.routes {
		handlers := make([]gin.HandlerFunc, 0)
		handlers = append(handlers, a.getRequestInstanceMiddleware)
		for _, mw := range route.Middlewares {
			middleware := a.middlewareToGinMiddleware(mw)
			handlers = append(handlers, middleware)
		}
		handler := a.handlerToGinHandler(route.Handler)
		handlers = append(handlers, handler)
		a.addRoute(route, handlers, api)
	}
}

func (a *Api) addRoute(route *types.Route, handlers []gin.HandlerFunc, api gin.IRoutes) {
	switch route.Method {
	case "GET":
		api.GET(route.Path, handlers...)
	case "POST":
		api.POST(route.Path, handlers...)
	case "PUT":
		api.PUT(route.Path, handlers...)
	case "DELETE":
		api.DELETE(route.Path, handlers...)
	case "PATCH":
		api.PATCH(route.Path, handlers...)
	}
}

func (a *Api) Run() error {
	ServerAddress := a.config.GetApiConfig().ServerHost + ":" + a.config.GetApiConfig().Port
	return a.engine.Run(ServerAddress)
}

func (a *Api) AuthenticateMiddleware(request *types.ApiRequest) *types.ApiResponse {
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		a.logger.Debug(apiInterfaces.MsgAuthHeaderMissing)
		return UnauthorizedResponse(apiInterfaces.MsgAuthHeaderMissing)
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	user, err := a.authManager.Authenticate(token)
	request.User = user
	if err != nil {
		a.logger.Debug("Authentication failed", "err", err)
		return UnauthorizedResponse("Invalid token")
	}
	return nil
}

func (a *Api) getRequestInstanceMiddleware(ctx *gin.Context) {
	apiRequest := a.ginCtxToApiRequest(ctx)
	ctx.Set(contextRequestKey, apiRequest)
	ctx.Next()
}

func (a *Api) handlerToGinHandler(handler func(request *types.ApiRequest) *types.ApiResponse) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiRequest := c.MustGet(contextRequestKey).(*types.ApiRequest)
		response := handler(apiRequest)
		for key, value := range response.Headers {
			c.Header(key, value)
		}
		c.JSON(response.StatusCode, response.Body)
	}
}

func (a *Api) middlewareToGinMiddleware(middleware types.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiRequest := c.MustGet(contextRequestKey).(*types.ApiRequest)
		response := middleware(apiRequest)
		if response != nil && !response.Body.Success {
			for key, value := range response.Headers {
				c.Header(key, value)
			}
			c.AbortWithStatusJSON(response.StatusCode, response.Body)
		}
		c.Next()
	}
}

func (a *Api) ginCtxToApiRequest(c *gin.Context) *types.ApiRequest {
	headers := getHeadersFromGinContext(c)
	queryParams := getQueryParamsFromGinContext(c)
	pathParams := make(map[string]string)
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Failed to read request body", "error", err)
	}
	return &types.ApiRequest{
		Method:      c.Request.Method,
		URL:         c.Request.URL.String(),
		Headers:     headers,
		QueryParams: queryParams,
		PathParams:  pathParams,
		Body:        string(bodyBytes),
		ClientIP:    c.ClientIP(),
	}
}

func (a *Api) healthCheck(request *types.ApiRequest) *types.ApiResponse {
	return &types.ApiResponse{
		StatusCode: http.StatusOK,
		Body: &types.ApiResponseBody{
			Success: true,
			Data: map[string]any{
				"status":  "ok",
				"message": "Wiki Scraper Newsletter API is running",
				"version": "1.0.0",
			},
		},
	}
}

func (a *Api) LoadModule(module apiInterfaces.ApiModuleInterface) {
	routes := module.GetRoutes()
	a.RegisterRoutes(routes, module.GetRoutePrefix())
}

func (a *Api) GetLogger() interfaces.LoggerInterface {
	return a.logger
}

func (a *Api) GetAuthManager() authType.AuthManagerInterface {
	return a.authManager
}

func (a *Api) recoverFromPanicMiddleware(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			a.logger.Error("Panic recovered in API middleware/handler", "error", r)
			a.logger.Error("Stack trace:\n" + stack)
			apiErr := InternalErrorResponse()
			ctx.AbortWithStatusJSON(apiErr.StatusCode, apiErr.Body)
		}
	}()
	ctx.Next()
}

func makeMiddlewareFromPolicy[T interfaces.ModelInterface](policyFunc func(user *model.User, resource T) error, resourceRepo pkgRepo.RepositoryInterface[T]) types.Middleware {
	return func(request *types.ApiRequest) *types.ApiResponse {
		user := request.User
		stringId := request.PathParams["id"]
		var resource T
		if stringId != "" {
			id, err := strconv.ParseUint(stringId, 10, 32)
			if err != nil {
				return BadRequestResponseWithMsg("Invalid filter ID")
			}
			resource, err = resourceRepo.GetByID(uint(id))
			if err != nil || any(resource) == nil {
				return NotFoundResponse("Filter not found")
			}
		}
		if err := policyFunc(user, resource); err != nil {
			return ForbiddenResponseWithMsg(err.Error())
		}
		return nil
	}
}
