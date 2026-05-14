package testutils

import "example-wikipedia-scraper/internal/model"

type MockUserRepository struct {
	MockGenericRepository[*model.User]
}

func (r *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := r.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := r.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *MockUserRepository) GetByEmailVerificationToken(token string) (*model.User, error) {
	args := r.Called(token)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *MockUserRepository) UpdateLastLogin(user *model.User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *MockUserRepository) FindByLogin(login string) (*model.User, error) {
	args := r.Called(login)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *MockUserRepository) GetByPasswordResetToken(token string) (*model.User, error) {
	args := r.Called(token)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *MockUserRepository) UpdateLastLogout(user *model.User) error {
	args := r.Called(user)
	return args.Error(0)
}
