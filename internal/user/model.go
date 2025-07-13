package user

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	MongoIDField       = "_id"
	ErrHashingPassword = "error hashing password"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string    `bson:"email,omitempty" json:"email,omitempty"`
	Password  string    `bson:"password,omitempty" json:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (u *User) HashPassword() error {
	if u.Password == "" {
		return errors.New(ErrHashingPassword)
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashBytes)

	return nil
}

func (u *User) ToBSONMap() (bson.M, error) {
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
