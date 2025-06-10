package libs

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte(os.Getenv("ACCESS_SECRET"))  // Use env in production!
	refreshSecret = []byte(os.Getenv("REFRESH_SECRET")) // Use env in production!
)

type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
	IV       string `json:"iv"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user UserClaims) (string, error) {
	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(user UserClaims) (string, error) {
	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(refreshSecret)
}

func GenerateRefreshTokenRememberMe(user UserClaims) (string, error) {
	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(refreshSecret)
}

func ParseAccessToken(tokenStr string) (*UserClaims, error) {
	return parseToken(tokenStr, accessSecret)
}

func ParseRefreshToken(tokenStr string) (*UserClaims, error) {
	return parseToken(tokenStr, refreshSecret)
}

func parseToken(tokenStr string, secret []byte) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
