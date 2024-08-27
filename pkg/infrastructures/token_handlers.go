package infrastructures

import (
	"errors"
	"loan-management/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.StandardClaims
}

func GenerateVerificationToken(userID, email, tokenType string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	config, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	return token.SignedString([]byte(config.Jwt.Secret))
}

func ValidateVerificationToken(tokenString, expectedEmail, expectedID, tokenType string) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.Secret), nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Email != expectedEmail || claims.UserID != expectedID || claims.Type != tokenType {
			return errors.New("invalid token")
		}
		if claims.ExpiresAt < time.Now().Unix() {
			return errors.New("token has expired")
		}
		return nil
	}

	return errors.New("invalid token")
}
