package main

import (
	"github.com/brandoyts/go-token-auth/internal/user"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type handler struct {
	userHandler *user.Handler
}

type appDependency struct {
	db      *mongo.Database
	redis   *redis.Client
	handler *handler
}
