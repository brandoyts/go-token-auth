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

func (rtr *RefreshTokenRepository) FindOne(ctx context.Context, model auth.RefreshToken) (*auth.RefreshToken, error) {
	filter, err := model.ToBSONMap()
	if err != nil {
		return nil, err
	}

	var result auth.RefreshToken

	err = rtr.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (rtr *RefreshTokenRepository) Update(ctx context.Context, id string, model auth.RefreshToken) error {
	filter, err := model.ToBSONMap()
	if err != nil {
		return err
	}

	_, err = rtr.collection.UpdateByID(ctx, id, bson.M{
		"$set": filter,
	})
	if err != nil {
		return err
	}

	return nil
}
