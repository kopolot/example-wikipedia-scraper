package api

import (
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/testutils"
	types "example-wikipedia-scraper/internal/types/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewApi_DebugMode(t *testing.T) {
	logger := &testutils.MockLogger{}
	config := &config.Config{ApiConfig: &config.ApiConfig{Debug: true, ServerHost: "127.0.0.1", Port: "8080"}}
	api := NewApi(config, logger, &testutils.MockAuthManager{})
	assert.NotNil(t, api)
	assert.Equal(t, config, api.config)
	assert.Equal(t, logger, api.logger)
	assert.NotNil(t, api.authManager)
	assert.NotNil(t, api.engine)
	assert.NotNil(t, api.routes)
}

func TestApi_HealthCheck(t *testing.T) {
	api := &Api{}
	req := &types.ApiRequest{}
	resp := api.healthCheck(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, resp.Body.Success)
	data, ok := resp.Body.Data.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "ok", data["status"])
}

func TestApi_GinCtxToApiRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/test?foo=bar", nil)
	c.Params = gin.Params{{Key: "id", Value: "123"}}
	api := &Api{}
	req := api.ginCtxToApiRequest(c)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/api/test?foo=bar", req.URL)
	assert.Equal(t, "123", req.PathParams["id"])
}

func TestApi_AuthenticateMiddleware_NoAuthHeader(t *testing.T) {
	logger := &testutils.MockLogger{}
	authManager := &testutils.MockAuthManager{}
	api := &Api{logger: logger, authManager: authManager}
	req := &types.ApiRequest{Headers: map[string]string{}}
	resp := api.AuthenticateMiddleware(req)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}

func TestApi_AuthenticateMiddleware_InvalidToken(t *testing.T) {
	logger := &testutils.MockLogger{}
	authManager := &testutils.MockAuthManager{}
	authManager.On("Authenticate", "invalid").Return((*model.User)(nil), assert.AnError)
	api := &Api{logger: logger, authManager: authManager}
	req := &types.ApiRequest{Headers: map[string]string{"Authorization": "Bearer invalid"}}
	resp := api.AuthenticateMiddleware(req)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}

func TestApi_AuthenticateMiddleware_ValidToken(t *testing.T) {
	logger := &testutils.MockLogger{}
	user := &model.User{Email: "test@example.com"}
	authManager := &testutils.MockAuthManager{}
	authManager.On("Authenticate", "valid").Return(user, nil)
	api := &Api{logger: logger, authManager: authManager}
	req := &types.ApiRequest{Headers: map[string]string{"Authorization": "Bearer valid"}}
	resp := api.AuthenticateMiddleware(req)
	assert.Nil(t, resp)
	assert.Equal(t, user, req.User)
}

func TestApi_RegisterRoutes(t *testing.T) {
	api := &Api{routes: make(map[string]*types.Route)}
	route := &types.Route{Method: "get", Path: "/test", Handler: func(r *types.ApiRequest) *types.ApiResponse { return nil }}
	api.RegisterRoutes([]*types.Route{route}, "")
	assert.Contains(t, api.routes, "/test_GET")
}

func TestApi_LoadModule(t *testing.T) {
	api := &Api{routes: make(map[string]*types.Route)}
	module := &testutils.MockApiModule{}
	module.On("GetRoutes").Return([]*types.Route{
		{
			Method:  "GET",
			Path:    "/mod",
			Handler: func(r *types.ApiRequest) *types.ApiResponse { return nil },
		},
	})
	module.On("GetRoutePrefix").Return("")
	api.LoadModule(module)
	assert.Contains(t, api.routes, "/mod_GET")
}

func TestApi_GetLogger(t *testing.T) {
	logger := &testutils.MockLogger{}
	api := &Api{logger: logger}
	assert.Equal(t, logger, api.GetLogger())
}

func TestApi_GetAuthManager(t *testing.T) {
	am := &testutils.MockAuthManager{}
	api := &Api{authManager: am}
	assert.Equal(t, am, api.GetAuthManager())
}
