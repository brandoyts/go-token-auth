package auth

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

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
		return nil, errors.New(invalidCredentials)
	}

	// generate access token
	accessToken, err := s.token.Generate(user.ID, os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshToken, err := s.token.Generate(user.ID, os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	ttlDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	ttl := time.Now().Add(ttlDuration)

	// store refresh token to db
	_, err = s.refreshTokenRepository.Create(ctx, RefreshToken{
		UserID:    user.ID,
		IPAddress: in.IPAddress,
		TokenHash: refreshToken,
		Revoked:   false,
		TTL:       ttl,
	})
	if err != nil {
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
		return nil, errors.New("suspicious")
	}

	userId := oldRefreshToken.UserID

	// mark old refresh token as revoked
	oldRefreshToken.Revoked = true
	err = s.refreshTokenRepository.Update(ctx, oldRefreshToken.ID, *oldRefreshToken)
	if err != nil {
		return nil, err
	}

	// generate new access token
	newAccessToken, err := s.token.Generate(userId, os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	// generate new refresh token
	newRefreshToken, err := s.token.Generate(userId, os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	ttlDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	ttl := time.Now().Add(ttlDuration)

	// save new refresh token to db
	_, err = s.refreshTokenRepository.Create(ctx, RefreshToken{
		UserID:    userId,
		IPAddress: in.IPAddress,
		Revoked:   false,
		TokenHash: newRefreshToken,
		TTL:       ttl,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	oldAccessToken, err := s.token.GetClaims(in.AccessToken)
	if err != nil {
		return nil, err
	}

	oldAccessTokenExpiration, err := oldAccessToken.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	cacheTtl := computeCacheExpiration(oldAccessTokenExpiration.Time)

	cacheTtlDuration := time.Until(cacheTtl).String()

	// blacklist old access token
	err = s.cache.Set(in.AccessToken, in.AccessToken, cacheTtlDuration)
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

	oldAccessToken, err := s.token.GetClaims(accessToken)
	if err != nil {
		return err
	}

	oldAccessTokenExpiration, err := oldAccessToken.GetExpirationTime()
	if err != nil {
		return err
	}

	cacheTtl := computeCacheExpiration(oldAccessTokenExpiration.Time)

	cacheTtlDuration := time.Until(cacheTtl).String()

	// blacklist access token
	return s.cache.Set(accessToken, accessToken, cacheTtlDuration)
}

func computeCacheExpiration(tokenTtl time.Time) time.Time {
	buffer := 30 * time.Second
	exp := tokenTtl.Add(-buffer)

	if exp.Before(time.Now()) {
		return time.Now().Add(buffer)
	}

	return exp
}
