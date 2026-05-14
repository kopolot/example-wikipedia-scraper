package service

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/logger"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/service"
	"example-wikipedia-scraper/internal/service/mailer"
	"example-wikipedia-scraper/internal/utilities"
	"example-wikipedia-scraper/test/integration"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	password          = "password"
	hashedPassword, _ = utilities.HashPassword(password)
	testUser          = &model.User{
		Email:    "test@example.com",
		Password: hashedPassword,
		Username: "test_user_1111",
	}
)

func setupTestUser(t *testing.T, userRepo *repository.UserRepository) *model.User {
	err := userRepo.Create(testUser)
	require.NoError(t, err)
	return testUser
}

func TestCreateUser_Success(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	logger := logger.GetLogger()
	mailerService := mailer.NewMailer(cfg.GetMailerConfig(), logger)
	userService := service.NewUserService(userRepo, mailerService, cfg)
	dto := dto.CreateUserDTO{
		Email:          testUser.Email,
		Password:       password,
		RepeatPassword: password,
		Username:       testUser.Username,
	}
	createdUser, err := userService.CreateUser(dto)
	assert.NoError(t, err)
	assert.Equal(t, dto.Email, createdUser.Email)
	assert.Equal(t, dto.Username, createdUser.Username)
	assert.True(t, utilities.CheckPasswordHash(password, createdUser.Password))
	userRepo.Delete(createdUser.ID)
}

func TestUpdateUser_Success(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	testUser := setupTestUser(t, userRepo)
	logger := logger.GetLogger()
	mailerService := mailer.NewMailer(cfg.GetMailerConfig(), logger)
	userService := service.NewUserService(userRepo, mailerService, cfg)
	newUsername := "updated_username_2222"
	newPassword := "new_password_2222"
	updateDTO := dto.UpdateUserDTO{
		Username: newUsername,
		Password: newPassword,
	}
	updatedUser, err := userService.UpdateUser(testUser, updateDTO)
	assert.NoError(t, err)
	assert.Equal(t, newUsername, updatedUser.Username)
	assert.True(t, utilities.CheckPasswordHash(newPassword, updatedUser.Password))
}

func TestGenerateVerificationToken_IsUUID(t *testing.T) {
	cfg, err := integration.InitTest()
	t.Cleanup(func() { integration.CleanupDB() })
	require.NoError(t, err)
	userRepo := repository.NewUserRepository()
	logger := logger.GetLogger()
	mailerService := mailer.NewMailer(cfg.GetMailerConfig(), logger)
	userService := service.NewUserService(userRepo, mailerService, cfg)
	token := userService.GenerateVerificationToken()
	fmt.Println("Generated token:", token)
	assert.NoError(t, uuid.Validate(token), "Generated token is not a valid UUID")
}
