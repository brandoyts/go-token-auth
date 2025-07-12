package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongodb(database string, uri string, credentials options.Credential) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), &readpref.ReadPref{})
	if err != nil {
		return nil, err
	}

	return client.Database(database), nil
}
