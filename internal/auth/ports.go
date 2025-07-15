package auth

import "context"

type DBInterface interface {
	Create(ctx context.Context, model RefreshToken) (string, error)
	FindOne(ctx context.Context, model RefreshToken) (*RefreshToken, error)
	Update(ctx context.Context, id string, model RefreshToken) error
}

type Cache interface {
	Set(key string, value string, ttl string) error
	Get(Key string) (string, error)
	Delete(key string) error
}

type Hash interface {
	Generate(value string) (string, error)
	Compare(hashed string, value string) error
}

type Token interface {
	Generate(tokenType string, id string) (string, error)
	Verify(tokenString string) error
}
