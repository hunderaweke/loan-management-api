package infrastructures

import (
	"fmt"
	"loan-management/config"
	"loan-management/internal/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	jwt.StandardClaims
	UserID  string
	Email   string
	IsAdmin bool
}

func GenerateJWTToken(user domain.User, t time.Duration) (string, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	claims := UserClaims{
		UserID:         user.ID,
		Email:          user.Email,
		IsAdmin:        user.IsAdmin,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(t).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Jwt.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (*UserClaims, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, fmt.Errorf("token has expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
