package mongodb

import (
	"context"

	"github.com/brandoyts/go-token-auth/internal/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const collection = "users"

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{collection: db.Collection(collection)}
}

func (ur *UserRepository) All(ctx context.Context) ([]user.User, error) {
	return []user.User{}, nil
}

func (ur *UserRepository) FindById(ctx context.Context, id string) (*user.User, error) {
	return nil, nil
}
func (ur *UserRepository) FindOne(ctx context.Context, user user.User) (*user.User, error) {
	return nil, nil
}
func (ur *UserRepository) Delete(ctx context.Context, id string) error {
	return nil
}
