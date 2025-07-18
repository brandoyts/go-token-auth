package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/brandoyts/go-token-auth/internal/auth"
	"github.com/brandoyts/go-token-auth/internal/infrastructure/hash"
	"github.com/brandoyts/go-token-auth/internal/infrastructure/jwtAuth"
	"github.com/brandoyts/go-token-auth/internal/infrastructure/mongodb"
	"github.com/brandoyts/go-token-auth/internal/infrastructure/redisClient"
	"github.com/brandoyts/go-token-auth/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func loadDependencies() *appDependency {

	// mongodb
	db, err := mongodb.NewMongodb(os.Getenv("MONGO_DATABASE_NAME"), os.Getenv("MONGO_URI"), options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})
	if err != nil {
		log.Fatal("❌ can't connect to mongodb", err)
	}

	fmt.Println("✅ successfully connected to mongodb")

	var tlsConfig *tls.Config

	if os.Getenv("REDIS_SSL") == "" {
		tlsConfig = nil
	}

	// redis
	redisClient := redisClient.NewRedisClient(&redis.Options{
		Addr:      os.Getenv("REDIS_ADDRESS"),
		Username:  os.Getenv("REDIS_USERNAME"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		TLSConfig: tlsConfig,
	})

	fmt.Println("✅ successfully connected to redis")

	// jwt
	jwtProvider := jwtAuth.New(os.Getenv("JWT_SECRET"))

	// inject user module
	userRepository := mongodb.NewUserRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// inject auth module
	hash := hash.New()
	refreshTokenRepository := mongodb.NewRefreshTokenRepository(db)
	authService := auth.NewService(hash, userService, jwtProvider, refreshTokenRepository, redisClient)
	authHandler := auth.NewHandler(authService)

	fmt.Println("✅ dependencies are loaded successfully")

	return &appDependency{
		db:          db,
		redis:       redisClient,
		jwtProvider: jwtProvider,
		handler: &handler{
			userHandler: userHandler,
			authHandler: authHandler,
		},
	}

}

func main() {

	deps := loadDependencies()

	// gracefully close mongodb connection
	defer func() {
		err := deps.db.Client().Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("⚠️ mongodb connection has been disconnected")
	}()

	// gracefully close redis connection

	app := fiber.New()
	app.Use(logger.New(), recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	apiRouter := app.Group("/api/v1")

	// user router
	userRouter := apiRouter.Group("/users")
	userRouter.Post("/create", deps.handler.userHandler.CreateUser)
	userRouter.Post("/find", deps.handler.userHandler.FindUser)
	userRouter.Get("/:id", deps.handler.userHandler.FindUserById)

	// auth router
	authRouter := apiRouter.Group("/auth")
	authRouter.Post("/login", deps.handler.authHandler.Login)
	authRouter.Post("/refresh-token", deps.handler.authHandler.RefreshToken)
	authRouter.Post("/logout", authChecker(deps.redis, deps.jwtProvider), deps.handler.authHandler.Logout)

	fmt.Println("server is listening on port 8080")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	app.Listen(fmt.Sprintf(":%v", port))
}
