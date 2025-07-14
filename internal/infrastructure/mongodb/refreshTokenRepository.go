package mongodb

import (
	"context"

	"github.com/brandoyts/go-token-auth/internal/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		collection: db.Collection("refresh_tokens"),
	}
}

func (rtr *RefreshTokenRepository) Create(ctx context.Context, model auth.RefreshToken) (string, error) {
	model.ID = bson.NewObjectID().Hex()

	result, err := rtr.collection.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}

	insertedID := result.InsertedID.(string)

	return insertedID, nil
}
