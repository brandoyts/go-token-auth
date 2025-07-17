package jwtAuth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuth struct {
	secret []byte
}

const (
	defaultSecretKey = "default-secret-key"
	empty            = ""
)

// error related constants
const (
	ErrParseJwt   = "error parsing jwt"
	ErrInvalidJwt = "invalid jwt"
)

func New(secret string) *JwtAuth {
	if secret == "" {
		secret = defaultSecretKey
	}

	return &JwtAuth{secret: []byte(secret)}
}

func (ja *JwtAuth) Generate(id string, ttl string) (string, error) {
	expiry, err := time.ParseDuration(ttl)
	if err != nil {
		return empty, err
	}

	claims := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(ja.secret)
	if err != nil {
		return empty, err
	}

	return tokenString, nil
}

func (ja *JwtAuth) Verify(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrParseJwt)
		}
		return ja.secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New(ErrInvalidJwt)
	}

	return nil
}

func (ja *JwtAuth) GetClaims(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrParseJwt)
		}
		return ja.secret, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}
