package main

import (
	"github.com/brandoyts/go-token-auth/internal/auth"
	"github.com/brandoyts/go-token-auth/internal/infrastructure/jwtAuth"
	"github.com/brandoyts/go-token-auth/internal/infrastructure/redisClient"
	"github.com/brandoyts/go-token-auth/internal/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type handler struct {
	userHandler *user.Handler
	authHandler *auth.Handler
}

type appDependency struct {
	db          *mongo.Database
	redis       *redisClient.RedisClient
	jwtProvider *jwtAuth.JwtAuth
	handler     *handler
}
