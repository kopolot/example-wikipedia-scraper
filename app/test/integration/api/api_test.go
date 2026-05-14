package api

import (
	"encoding/json"
	"example-wikipedia-scraper/internal/api"
	"example-wikipedia-scraper/internal/auth"
	"example-wikipedia-scraper/internal/config"
	"example-wikipedia-scraper/internal/dto"
	apiInterfaces "example-wikipedia-scraper/internal/interfaces/api"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/service"
	"example-wikipedia-scraper/internal/service/mailer"
	apiTypes "example-wikipedia-scraper/internal/types/api"
	"example-wikipedia-scraper/test/integration"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userPassword = "securepassword"
	testUser     = &model.User{
		Email:    "test@example.com",
		Username: "test12312123",
	}
	appConfig config.ConfigInterface
)

func setupApiStruct(t *testing.T) *api.Api {
	cfg, err := integration.InitTest()
	appConfig = cfg
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	authManager := auth.NewAuthManager(
		cfg.GetApiConfig(),
		userRepo,
	)
	api := api.NewApi(
		cfg,
		logger.GetLogger(),
		authManager,
	)
	return api
}

func setupUserService(cfg config.ConfigInterface) *service.UserService {
	userService := service.NewUserService(
		repository.NewUserRepository(),
		mailer.NewMailer(cfg.GetMailerConfig(), logger.GetLogger()),
		cfg,
	)
	return userService
}

func createOrGetTestUser(t *testing.T) *model.User {
	t.Helper()
	userRepo := repository.NewUserRepository()
	existingUser, err := userRepo.GetByEmail(testUser.Email)
	if err == nil && existingUser != nil && existingUser.ID != 0 {
		return existingUser
	}
	userService := setupUserService(appConfig)
	createdUser, err := userService.CreateUser(dto.CreateUserDTO{
		Email:    testUser.Email,
		Username: testUser.Username,
		Password: userPassword,
	})
	require.NoError(t, err)
	return createdUser
}

func TestAuthenticateMiddleware_HeaderMissing(t *testing.T) {
	api := setupApiStruct(t)
	t.Cleanup(func() { integration.CleanupDB() })
	req := &apiTypes.ApiRequest{}
	resp := api.AuthenticateMiddleware(req)
	require.NotNil(t, resp)
	require.Equal(t, 401, resp.StatusCode)
	require.Equal(t, apiInterfaces.MsgAuthHeaderMissing, resp.Body.Errors[0].Message)
}

func TestAuthenticateMiddleware_TokenInvalid(t *testing.T) {
	api := setupApiStruct(t)
	t.Cleanup(func() { integration.CleanupDB() })
	req := &apiTypes.ApiRequest{
		Headers: map[string]string{
			"Authorization": "Bearer invalid_token",
		},
	}
	resp := api.AuthenticateMiddleware(req)
	require.NotNil(t, resp)
	require.Equal(t, 401, resp.StatusCode)
}

func TestAuthenticateMiddleware_Success(t *testing.T) {
	api := setupApiStruct(t)
	t.Cleanup(func() { integration.CleanupDB() })
	user := createOrGetTestUser(t)
	token, err := api.GetAuthManager().GenerateToken(*user, time.Second*10)
	require.NoError(t, err)
	req := &apiTypes.ApiRequest{
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	}
	resp := api.AuthenticateMiddleware(req)
	assert.Nil(t, resp)
}

func TestHealthCheckEndpoint(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	api := setupApiStruct(t)
	api.SetupRoutes()
	go func() {
		err = api.Run()
	}()
	time.Sleep(time.Second)
	assert.NoError(t, err)
	url := "http://localhost:" + cfg.GetApiConfig().Port + "/api/health"
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	jsonBody := &apiTypes.ApiResponseBody{}
	err = json.Unmarshal(body, jsonBody)
	assert.NoError(t, err)
	assert.True(t, jsonBody.Success)
}

func TestLoadModules(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	api := setupApiStruct(t)
	api.LoadModules()
	api.SetupRoutes()
	go func() {
		err = api.Run()
	}()
	time.Sleep(time.Second)
	assert.NoError(t, err)
	url := "http://localhost:" + cfg.GetApiConfig().Port + "/api/health"
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	jsonBody := &apiTypes.ApiResponseBody{}
	err = json.Unmarshal(body, jsonBody)
	assert.NoError(t, err)
	assert.True(t, jsonBody.Success)
}

func TestPanicInRoute(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	api := setupApiStruct(t)
	api.RegisterRoutes([]*apiTypes.Route{
		{
			Method: "GET",
			Path:   "/testpanic",
			Handler: func(r *apiTypes.ApiRequest) *apiTypes.ApiResponse {
				panic("test panic!")
			},
		},
	}, "")
	api.SetupRoutes()
	go func() {
		err = api.Run()
	}()
	time.Sleep(time.Second)
	assert.NoError(t, err)
	url := "http://localhost:" + cfg.GetApiConfig().Port + "/api/testpanic"
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	jsonBody := &apiTypes.ApiResponseBody{}
	err = json.Unmarshal(body, jsonBody)
	assert.NoError(t, err)
	assert.False(t, jsonBody.Success)
	assert.Equal(t, "Internal Server Error", jsonBody.Errors[0].Message)
}
