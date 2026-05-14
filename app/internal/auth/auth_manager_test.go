package auth

import (
	"example-wikipedia-scraper/internal/config"
	authInterfaces "example-wikipedia-scraper/internal/interfaces/auth"
	"example-wikipedia-scraper/internal/model"
	testRepo "example-wikipedia-scraper/internal/testutils/repository"
	"example-wikipedia-scraper/internal/types/auth"
	"example-wikipedia-scraper/internal/utilities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	apiConfig = &config.ApiConfig{
		JWT: &config.JwtConfig{
			Secret:         "test-secret",
			ExpirationTime: 3600,
		},
	}
	testUserId                 uint = 1
	testUserEmail                   = "test@example.com"
	testUsername                    = "testuser"
	testUserPassword                = "testpassword"
	hashedPassword, _               = utilities.HashPassword(testUserPassword)
	userSubscriptionExpiration      = time.Now().Add(time.Hour)
	testUserModel                   = &model.User{
		Model: model.Model{
			ID: testUserId,
		},
		Email:         testUserEmail,
		Username:      testUsername,
		Password:      hashedPassword,
		EmailVerified: true,
	}
)

func TestAuthenticate_Success(t *testing.T) {
	userRepo := &testRepo.MockUserRepository{}
	userRepo.MockGenericRepository.On("GetByID", uint(1)).Return(testUserModel, nil)
	userRepo.On("UpdateLastLogin", testUserModel).Return(nil)
	authManager := NewAuthManager(apiConfig, userRepo)
	authManager.parseJWT = func(secret []byte, token string) (*auth.JWTClaims, error) {
		return &auth.JWTClaims{UserID: testUserId, Email: testUserEmail, ExpiresAt: time.Now().Add(time.Minute).Unix()}, nil
	}
	model, err := authManager.Authenticate("test-token")
	assert.NoError(t, err)
	require.NotNil(t, model)
	assert.Equal(t, *testUserModel, *model)
	userRepo.AssertExpectations(t)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	userRepo := &testRepo.MockUserRepository{}
	authManager := NewAuthManager(apiConfig, userRepo)
	authManager.parseJWT = func(secret []byte, token string) (*auth.JWTClaims, error) {
		return nil, assert.AnError
	}
	user, err := authManager.Authenticate("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestAuthenticate_UserNotFound(t *testing.T) {
	userRepo := &testRepo.MockUserRepository{}
	userRepo.MockGenericRepository.On("GetByID", uint(1)).Return(&model.User{}, nil)
	authManager := NewAuthManager(apiConfig, userRepo)
	authManager.parseJWT = func(secret []byte, token string) (*auth.JWTClaims, error) {
		return &auth.JWTClaims{UserID: testUserId, Email: testUserEmail}, nil
	}
	user, err := authManager.Authenticate("test-token")
	assert.Error(t, err)
	assert.Nil(t, user)
	userRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	userRepo := &testRepo.MockUserRepository{}
	userRepo.On("FindByLogin", testUsername).Return(testUserModel, nil)
	userRepo.On("UpdateLastLogin", testUserModel).Return(nil)
	authManager := NewAuthManager(apiConfig, userRepo)
	token, err := authManager.Login(testUsername, testUserPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	userRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	userRepo := &testRepo.MockUserRepository{}
	userRepo.On("FindByLogin", testUsername).Return(&model.User{}, nil)
	authManager := NewAuthManager(apiConfig, userRepo)
	token, err := authManager.Login(testUsername, "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, authInterfaces.ErrInvalidCredentials, err)
	assert.Empty(t, token)
	userRepo.AssertExpectations(t)
}
