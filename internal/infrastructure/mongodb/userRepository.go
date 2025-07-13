package mongodb

import (
	"context"

	"github.com/brandoyts/go-token-auth/internal/user"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	filter := bson.D{{user.MongoIDField, id}}

	var result user.User

	err := ur.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}
func (ur *UserRepository) FindOne(ctx context.Context, usr user.User) (*user.User, error) {

	bsonMap, err := usr.ToBSONMap()
	if err != nil {
		return nil, err
	}

	var result user.User
	err = ur.collection.FindOne(ctx, bsonMap).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}
func (ur *UserRepository) Create(ctx context.Context, user user.User) (string, error) {
	user.ID = bson.NewObjectID().Hex()

	result, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	insertedID := result.InsertedID.(string)

	return insertedID, nil
}
func (ur *UserRepository) Delete(ctx context.Context, id string) error {
	return nil
}
