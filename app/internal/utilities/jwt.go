package utilities

import (
	"example-wikipedia-scraper/internal/types/auth"
	"strconv"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
	jwt.RegisteredClaims
}

func GenerateJWT(jwtSecret []byte, claims auth.JWTClaims) (string, error) {
	idString := strconv.Itoa(int(claims.UserID))
	jwtClaims := jwtClaims{
		Username:  claims.Username,
		Email:     claims.Email,
		Role:      claims.Role,
		CreatedAt: claims.CreatedAt,
		UpdatedAt: claims.UpdatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(claims.ExpiresAt) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   idString,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(jwtSecret []byte, tokenString string) (*auth.JWTClaims, error) {
	jwtClaims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	userID, _ := strconv.Atoi(jwtClaims.RegisteredClaims.Subject)
	claims := &auth.JWTClaims{
		UserID:    uint(userID),
		Email:     jwtClaims.Email,
		Role:      jwtClaims.Role,
		IssuedAt:  jwtClaims.RegisteredClaims.IssuedAt.Unix(),
		ExpiresAt: jwtClaims.RegisteredClaims.ExpiresAt.Unix(),
		Username:  jwtClaims.Username,
		CreatedAt: jwtClaims.CreatedAt,
		UpdatedAt: jwtClaims.UpdatedAt,
	}
	if token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
