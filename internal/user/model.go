package user

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string    `bson:"email,omitempty" json:"email,omitempty"`
	Password  string    `bson:"password,omitempty" json:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
