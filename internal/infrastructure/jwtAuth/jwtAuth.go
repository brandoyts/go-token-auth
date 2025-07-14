package jwtauth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuth struct {
	secret []byte
}

func New(secret string) *JwtAuth {
	if secret == "" {
		secret = "default-secret-key"
	}

	return &JwtAuth{secret: []byte(secret)}
}

func (ja *JwtAuth) Generate(tokenType string, id string) (string, error) {
	claims := jwt.MapClaims{
		"token_type": tokenType,
		"id":         id,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(ja.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ja *JwtAuth) Verify(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return ja.secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
