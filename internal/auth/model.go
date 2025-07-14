package auth

import "time"

type RefreshToken struct {
	ID        string    `bson:"_id,omitempty"`
	IsRevoked bool      `bson:"is_revoked"`
	Token     string    `bson:"token,omitempty"`
	TTL       time.Time `bson:"ttl,omitempty"`
}
