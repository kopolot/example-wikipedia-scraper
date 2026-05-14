package api

import (
	"encoding/json"
	"errors"
	"testing"

	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces/auth"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/testutils"
	apitypes "example-wikipedia-scraper/internal/types/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	testRepo "example-wikipedia-scraper/internal/testutils/repository"
	testService "example-wikipedia-scraper/internal/testutils/service"
)

func newUserApiModuleMocks() (*UserApiModule, *testRepo.MockUserRepository, *testutils.MockApi, *testutils.MockLogger, *testutils.MockAuthManager, *testService.MockUserService) {
	userRepo := &testRepo.MockUserRepository{}
	userService := &testService.MockUserService{}
	logger := &testutils.MockLogger{}
	authManager := &testutils.MockAuthManager{}
	api := &testutils.MockApi{Logger: logger, AuthManager: authManager}
	return NewUserApiModule(userRepo, userService, api), userRepo, api, logger, authManager, userService
}

func TestCreateUser_EmailConflict(t *testing.T) {
	module, userRepo, _, _, _, userService := newUserApiModuleMocks()
	req := dto.CreateUserDTO{Email: "test@test.com", Username: "user1", Password: "password", RepeatPassword: "password"}
	user := &model.User{}
	user.ID = 1
	userRepo.On("GetByEmail", req.Email).Return(user, nil)
	body, _ := json.Marshal(req)
	apiReq := &apitypes.ApiRequest{Body: string(body)}
	resp := module.createUser(apiReq)
	assert.Equal(t, 409, resp.StatusCode)
	assert.False(t, resp.Body.Success)
	userRepo.AssertExpectations(t)
	userService.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestCreateUser_UsernameConflict(t *testing.T) {
	module, userRepo, _, _, _, userService := newUserApiModuleMocks()
	req := dto.CreateUserDTO{Email: "test@test.com", Username: "user1", Password: "password", RepeatPassword: "password"}
	userRepo.On("GetByEmail", req.Email).Return(&model.User{}, nil)
	user := &model.User{}
	user.ID = 2
	userRepo.On("GetByUsername", req.Username).Return(user, nil)
	body, _ := json.Marshal(req)
	apiReq := &apitypes.ApiRequest{Body: string(body)}
	resp := module.createUser(apiReq)
	assert.Equal(t, 409, resp.StatusCode)
	assert.False(t, resp.Body.Success)
	userRepo.AssertExpectations(t)
	userService.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestCreateUser_Success(t *testing.T) {
	module, userRepo, _, _, _, userService := newUserApiModuleMocks()
	req := dto.CreateUserDTO{Email: "test@test.com", Username: "user1", Password: "password", RepeatPassword: "password"}
	userRepo.On("GetByEmail", req.Email).Return(&model.User{}, nil)
	userRepo.On("GetByUsername", req.Username).Return(&model.User{}, nil)
	user := &model.User{Email: req.Email, Username: req.Username}
	user.ID = 3
	userService.On("CreateUser", req).Return(user, nil)
	body, _ := json.Marshal(req)
	apiReq := &apitypes.ApiRequest{Body: string(body)}
	resp := module.createUser(apiReq)
	assert.Equal(t, 201, resp.StatusCode)
	assert.True(t, resp.Body.Success)
	assert.Equal(t, uint(3), resp.Body.Data.(*model.User).ID)
	userRepo.AssertExpectations(t)
	userService.AssertExpectations(t)
}

func TestCreateUser_InvalidUser(t *testing.T) {
	module, _, _, _, _, userService := newUserApiModuleMocks()
	apiReq := &apitypes.ApiRequest{Body: "{\"email\":\"test\",\"username\":\"user1\", \"password\":\"pass\",\"repeat_password\":\"pass\"}"}
	resp := module.createUser(apiReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.False(t, resp.Body.Success)
	userService.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestCreateUser_InvalidBody(t *testing.T) {
	module, _, _, _, _, userService := newUserApiModuleMocks()
	apiReq := &apitypes.ApiRequest{Body: "{invalid json}"}
	resp := module.createUser(apiReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.False(t, resp.Body.Success)
	userService.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestLoginUser_InvalidBody(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	apiReq := &apitypes.ApiRequest{Body: "{invalid json}"}
	resp := module.loginUser(apiReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	module, _, _, _, authManager, _ := newUserApiModuleMocks()
	req := dto.LoginUserDTO{Login: "user", Password: "badpass"}
	body, _ := json.Marshal(req)
	apiReq := &apitypes.ApiRequest{Body: string(body)}
	authManager.On("Login", req.Login, req.Password).Return("", auth.ErrInvalidCredentials)
	resp := module.loginUser(apiReq)
	assert.Equal(t, 401, resp.StatusCode)
	assert.False(t, resp.Body.Success)
	authManager.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	module, _, _, _, authManager, _ := newUserApiModuleMocks()
	req := dto.LoginUserDTO{Login: "user", Password: "goodpass"}
	body, _ := json.Marshal(req)
	apiReq := &apitypes.ApiRequest{Body: string(body)}
	authManager.On("Login", req.Login, req.Password).Return("token123", nil)
	resp := module.loginUser(apiReq)
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, resp.Body.Success)
	assert.Equal(t, "token123", resp.Body.Data.(map[string]string)["token"])
	authManager.AssertExpectations(t)
}

func TestGetUser_InvalidID(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	apiReq := &apitypes.ApiRequest{PathParams: map[string]string{"id": "abc"}}
	resp := module.getUser(apiReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}

func TestGetUser_NotFound(t *testing.T) {
	module, userRepo, _, _, _, _ := newUserApiModuleMocks()
	userRepo.On("GetByID", uint(99)).Return(&model.User{}, nil)
	apiReq := &apitypes.ApiRequest{PathParams: map[string]string{"id": "99"}}
	resp := module.getUser(apiReq)
	assert.Equal(t, 404, resp.StatusCode)
	assert.False(t, resp.Body.Success)
	userRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	module, userRepo, _, _, _, _ := newUserApiModuleMocks()
	user := &model.User{Email: "a@b.com"}
	user.ID = 1
	userRepo.On("GetByID", uint(1)).Return(user, nil)
	apiReq := &apitypes.ApiRequest{PathParams: map[string]string{"id": "1"}}
	resp := module.getUser(apiReq)
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, resp.Body.Success)
	assert.Equal(t, "a@b.com", resp.Body.Data.(*model.User).Email)
	userRepo.AssertExpectations(t)
}

func TestDeleteUser_InvalidID(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	apiReq := &apitypes.ApiRequest{PathParams: map[string]string{"id": "abc"}}
	resp := module.deleteUser(apiReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}

func TestDeleteUser_Error(t *testing.T) {
	module, userRepo, _, _, _, _ := newUserApiModuleMocks()
	userRepo.On("Disable", uint(5)).Return(errors.New("fail"))
	apiReq := &apitypes.ApiRequest{PathParams: map[string]string{"id": "5"}}
	resp := module.deleteUser(apiReq)
	assert.NotEqual(t, 204, resp.StatusCode)
	assert.False(t, resp.Body.Success)
	userRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	module, userRepo, _, _, _, _ := newUserApiModuleMocks()
	userRepo.On("Disable", uint(5)).Return(nil)
	apiReq := &apitypes.ApiRequest{PathParams: map[string]string{"id": "5"}}
	resp := module.deleteUser(apiReq)
	assert.Equal(t, 204, resp.StatusCode)
	assert.True(t, resp.Body.Success)
	userRepo.AssertExpectations(t)
}

func TestAuthorizeUserMiddleware_Unauthorized(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	apiReq := &apitypes.ApiRequest{User: nil, PathParams: map[string]string{"id": "1"}}
	resp := module.AuthorizeUserMiddleware(apiReq)
	assert.Equal(t, 401, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}

func TestAuthorizeUserMiddleware_Admin(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	user := &model.User{Role: model.RoleAdmin}
	user.ID = 1
	apiReq := &apitypes.ApiRequest{User: user, PathParams: map[string]string{"id": "1"}}
	resp := module.AuthorizeUserMiddleware(apiReq)
	assert.Nil(t, resp)
}

func TestAuthorizeUserMiddleware_InvalidID(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	user := &model.User{Role: model.RoleUser}
	user.ID = 2
	apiReq := &apitypes.ApiRequest{User: user, PathParams: map[string]string{"id": "abc"}}
	resp := module.AuthorizeUserMiddleware(apiReq)
	assert.Equal(t, 400, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}

func TestAuthorizeUserMiddleware_Owner(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	user := &model.User{Role: model.RoleUser}
	user.ID = 2
	apiReq := &apitypes.ApiRequest{User: user, PathParams: map[string]string{"id": "2"}}
	resp := module.AuthorizeUserMiddleware(apiReq)
	assert.Nil(t, resp)
}

func TestAuthorizeUserMiddleware_Forbidden(t *testing.T) {
	module, _, _, _, _, _ := newUserApiModuleMocks()
	user := &model.User{Role: model.RoleUser}
	user.ID = 2
	apiReq := &apitypes.ApiRequest{User: user, PathParams: map[string]string{"id": "3"}}
	resp := module.AuthorizeUserMiddleware(apiReq)
	assert.Equal(t, 403, resp.StatusCode)
	assert.False(t, resp.Body.Success)
}
