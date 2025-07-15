package auth

import (
	"context"
	"errors"
	"log"

	"github.com/brandoyts/go-token-auth/internal/shared"
	"github.com/brandoyts/go-token-auth/internal/user"
)

const invalidCredentials = "invalid credentials"

type Service struct {
	userService            *user.Service
	hash                   Hash
	token                  Token
	refreshTokenRepository DBInterface
	cache                  Cache
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}

type LoginInput struct {
	IPAddress string
	Email     string
	Password  string
}

type RefreshTokenInput struct {
	IPAddress    string
	RefreshToken string
	AccessToken  string
}

func NewService(hash Hash, userService *user.Service, token Token, refreshTokenRepository DBInterface, cache Cache) *Service {
	return &Service{
		hash:                   hash,
		userService:            userService,
		token:                  token,
		refreshTokenRepository: refreshTokenRepository,
		cache:                  cache,
	}
}

func (s *Service) Login(ctx context.Context, in LoginInput) (*AuthToken, error) {
	// query user by email
	user, err := s.userService.FindUser(ctx, user.User{Email: in.Email})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(invalidCredentials)
	}

	// compare hashed password against login password
	err = s.hash.Compare(user.Password, in.Password)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New(invalidCredentials)
	}

	// generate access token
	accessToken, err := s.token.Generate(shared.TOKEN_TYPE_ACCESS, user.ID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// generate refresh token
	refreshToken, err := s.token.Generate(shared.TOKEN_TYPE_REFRESH, "")
	if err != nil {
		return nil, err
	}

	// store refresh token to db
	_, err = s.refreshTokenRepository.Create(ctx, RefreshToken{
		UserID:    user.ID,
		IPAddress: in.IPAddress,
		TokenHash: refreshToken,
		Revoked:   false,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, in RefreshTokenInput) (*AuthToken, error) {
	// query refresh token from db
	oldRefreshToken, err := s.refreshTokenRepository.FindOne(ctx, RefreshToken{TokenHash: in.RefreshToken})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if oldRefreshToken == nil {
		return nil, errors.New("token does not exists")
	}

	// check if token is already revoked
	if oldRefreshToken.Revoked {
		return nil, errors.New("token is already revoked")
	}

	// check request ip address if it match to old refresh token
	if in.IPAddress != oldRefreshToken.IPAddress {
		return nil, errors.New("suspicious!")
	}

	// mark old refresh token as revoked
	oldRefreshToken.Revoked = true
	err = s.refreshTokenRepository.Update(ctx, oldRefreshToken.ID, *oldRefreshToken)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// generate new access token
	newAccessToken, err := s.token.Generate(shared.TOKEN_TYPE_ACCESS, "")
	if err != nil {
		return nil, err
	}

	// generate new refresh token
	newRefreshToken, err := s.token.Generate(shared.TOKEN_TYPE_REFRESH, "")
	if err != nil {
		return nil, err
	}

	// save new refresh token to db
	_, err = s.refreshTokenRepository.Create(ctx, RefreshToken{
		UserID:    oldRefreshToken.UserID,
		IPAddress: in.IPAddress,
		Revoked:   false,
		TokenHash: newRefreshToken,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// blacklist old access token
	err = s.cache.Set(in.AccessToken, in.AccessToken, "1h")
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	// validate refresh token
	err := s.token.Verify(refreshToken)
	if err != nil {
		return nil
	}

	// check refresh token if it exist in db
	oldRefreshToken, err := s.refreshTokenRepository.FindOne(ctx, RefreshToken{
		TokenHash: refreshToken,
	})
	if err != nil {
		return nil
	}
	if oldRefreshToken == nil {
		return errors.New("token does not exists")
	}

	// mark refresh token as revoked
	err = s.refreshTokenRepository.Update(ctx, oldRefreshToken.ID, RefreshToken{
		Revoked: true,
	})
	if err != nil {
		return nil
	}

	// blacklist access token
	return s.cache.Set(accessToken, accessToken, "1h")
}
