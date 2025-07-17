package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/brandoyts/go-token-auth/internal/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {

	ttl := os.Getenv("REFRESH_TOKEN_TTL")
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		log.Fatal(err)
	}

	durationInSeconds := int32(duration.Seconds())

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "ttl", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(durationInSeconds),
	}

	collection := db.Collection("refresh_tokens")

	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

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
