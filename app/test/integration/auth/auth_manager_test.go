package auth

import (
	"example-wikipedia-scraper/internal/auth"
	"example-wikipedia-scraper/internal/config"
	authInterfaces "example-wikipedia-scraper/internal/interfaces/auth"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/utilities"
	"example-wikipedia-scraper/test/integration"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	password          = "password"
	hashedPassword, _ = utilities.HashPassword(password)
	testUser          = &model.User{
		Email:         "test@example.com",
		Password:      hashedPassword,
		Username:      "test_user_1111",
		EmailVerified: true,
	}
	appConfig config.ConfigInterface
)

func createOrGetTestUser(t *testing.T, userRepo *repository.UserRepository) *model.User {
	err := userRepo.Create(testUser)
	require.NoError(t, err)
	return testUser
}

func TestGenerateToken_Success(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	authManager := auth.NewAuthManager(
		cfg.GetApiConfig(),
		userRepo,
	)
	token, err := authManager.GenerateToken(*testUser, time.Duration(cfg.GetApiConfig().JWT.ExpirationTime)*time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_Success(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	user := createOrGetTestUser(t, userRepo)

	authManager := auth.NewAuthManager(
		cfg.GetApiConfig(),
		userRepo,
	)
	token, err := authManager.Login(user.Email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	authManager := auth.NewAuthManager(
		cfg.GetApiConfig(),
		userRepo,
	)
	_, err = authManager.Login("invalid_email@example.com", "invalid_password")
	assert.Error(t, err)
	assert.Equal(t, authInterfaces.ErrInvalidCredentials, err)
}

func TestAuthenticate_Success(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	user := createOrGetTestUser(t, userRepo)
	authManager := auth.NewAuthManager(
		cfg.GetApiConfig(),
		userRepo,
	)
	token, err := authManager.GenerateToken(*user, time.Duration(cfg.GetApiConfig().JWT.ExpirationTime)*time.Second)
	require.NoError(t, err)
	authenticatedUser, err := authManager.Authenticate(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, authenticatedUser.ID)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	authManager := auth.NewAuthManager(
		cfg.GetApiConfig(),
		userRepo,
	)
	_, err = authManager.Authenticate("invalid_token")
	assert.Error(t, err)
}
