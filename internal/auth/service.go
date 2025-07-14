package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/brandoyts/go-token-auth/internal/user"
)

const invalidCredentials = "invalid credentials"

type Service struct {
	userService            *user.Service
	hash                   Hash
	token                  Token
	refreshTokenRepository DBInterface
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}

func NewService(hash Hash, userService *user.Service, token Token, refreshTokenRepository DBInterface) *Service {
	return &Service{
		hash:                   hash,
		userService:            userService,
		token:                  token,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (s *Service) Login(ctx context.Context, usr user.User) (*AuthToken, error) {
	result, err := s.userService.FindUser(ctx, user.User{Email: usr.Email})
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New(invalidCredentials)
	}

	// compare hashed password against login password
	err = s.hash.Compare(result.Password, usr.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(invalidCredentials)
	}

	accessToken, err := s.token.Generate("access", result.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.token.Generate("refresh", result.ID)
	if err != nil {
		return nil, err
	}

	_, err = s.refreshTokenRepository.Create(ctx, RefreshToken{
		IsRevoked: false,
		Token:     refreshToken,
	})

	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
