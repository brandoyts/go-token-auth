package auth

import "context"

type DBInterface interface {
	Create(ctx context.Context, model RefreshToken) (string, error)
}

type Hash interface {
	Generate(value string) (string, error)
	Compare(hashed string, value string) error
}

type Token interface {
	Generate(tokenType string, id string) (string, error)
	Verify(tokenString string) error
}
