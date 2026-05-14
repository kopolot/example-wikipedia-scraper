package auth

import (
	"errors"
	"time"

	"example-wikipedia-scraper/internal/config"
	authInterfaces "example-wikipedia-scraper/internal/interfaces/auth"
	repoInterfaces "example-wikipedia-scraper/internal/interfaces/repository"
	"example-wikipedia-scraper/internal/model"
	"example-wikipedia-scraper/internal/types/auth"
	"example-wikipedia-scraper/internal/utilities"
)

type AuthManager struct {
	jwtSecret   []byte
	userRepo    repoInterfaces.UserRepositoryInterface
	apiConfig   *config.ApiConfig
	parseJWT    func(jwtSecret []byte, tokenString string) (*auth.JWTClaims, error)
	generateJWT func(jwtSecret []byte, claims auth.JWTClaims) (string, error)
}

func NewAuthManager(apiConfig *config.ApiConfig, userRepo repoInterfaces.UserRepositoryInterface) *AuthManager {
	return &AuthManager{
		jwtSecret:   []byte(apiConfig.JWT.Secret),
		userRepo:    userRepo,
		apiConfig:   apiConfig,
		parseJWT:    utilities.ParseJWT,
		generateJWT: utilities.GenerateJWT,
	}
}

func (a *AuthManager) Authenticate(token string) (*model.User, error) {
	claims, err := a.ParseToken(token)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	if !user.EmailVerified {
		return nil, authInterfaces.ErrUserNotVerified
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, authInterfaces.ErrTokenExpired
	}
	if user.LastLogoutAt != nil && user.LastLogoutAt.Unix() > claims.IssuedAt {
		return nil, authInterfaces.ErrTokenExpired
	}
	a.userRepo.UpdateLastLogin(user)
	return user, nil
}

func (a *AuthManager) GenerateToken(user model.User, expirationTime time.Duration) (string, error) {
	if expirationTime <= 0 {
		expirationTime = time.Duration(a.apiConfig.JWT.ExpirationTime) * time.Second
	}
	claims := auth.JWTClaims{
		Username:  user.Username,
		UserID:    user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		IssuedAt:  time.Now().Unix(),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
		ExpiresAt: time.Now().Add(
			expirationTime,
		).Unix(),
	}
	return a.generateJWT(a.jwtSecret, claims)
}

func (a *AuthManager) ParseToken(token string) (*auth.JWTClaims, error) {
	return a.parseJWT(a.jwtSecret, token)
}

func (a *AuthManager) Login(login string, password string) (string, error) {
	user, err := a.userRepo.FindByLogin(login)
	if err != nil {
		return "", authInterfaces.ErrInvalidCredentials
	}
	if user.ID == 0 || !utilities.CheckPasswordHash(password, user.Password) {
		return "", authInterfaces.ErrInvalidCredentials
	}
	if !user.EmailVerified {
		return "", authInterfaces.ErrUserNotVerified
	}
	token, err := a.GenerateToken(*user, time.Duration(a.apiConfig.JWT.ExpirationTime)*time.Second)
	if err != nil {
		return "", err
	}
	a.userRepo.UpdateLastLogin(user)
	return token, nil
}
