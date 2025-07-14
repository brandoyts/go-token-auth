package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshToken struct {
	ID        string    `bson:"_id,omitempty"`
	IsRevoked bool      `bson:"is_revoked"`
	Token     string    `bson:"token,omitempty"`
	TTL       time.Time `bson:"ttl,omitempty"`
}

func (u *RefreshToken) ToBSONMap() (bson.M, error) {
	var bsonMap bson.M

	filterBye, err := bson.Marshal(u)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(filterBye, &bsonMap)
	if err != nil {
		return nil, err
	}

	return bsonMap, nil
}
