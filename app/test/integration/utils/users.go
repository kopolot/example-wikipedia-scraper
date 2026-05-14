package utils

import (
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/model/repository"
	"example-wikipedia-scraper/internal/utilities"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func SetUpTestUser(t *testing.T) *model.User {
	userRepo := repository.NewUserRepository()
	user := &model.User{
		Email: "testuser@example.com",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)
	return user
}

func TearDownTestUser(t *testing.T, userID uint) {
	userRepo := repository.NewUserRepository()
	err := userRepo.Delete(userID)
	require.NoError(t, err)
}

func CreateTestUsers(t *testing.T, count int) []*model.User {
	userRepo := repository.NewUserRepository()
	var users []*model.User
	for i := range count {
		password := "------"
		hashedPassword, err := utilities.HashPassword(password)
		require.NoError(t, err)
		user := &model.User{
			Email:         "testuser" + strconv.Itoa(i) + "@example.com",
			Username:      fmt.Sprintf("testuser_%d", i),
			Password:      hashedPassword,
			EmailVerified: true,
		}
		users = append(users, user)
	}
	err := userRepo.CreateInBatches(users, uint(count))
	require.NoError(t, err)
	return users
}
