package user

import "context"

type DBInterface interface {
	All(ctx context.Context) ([]User, error)
	FindById(ctx context.Context, id string) (*User, error)
	FindOne(ctx context.Context, user User) (*User, error)
	Delete(ctx context.Context, id string) error
}
