package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/brandoyts/go-token-auth/internal/infrastructure/mongodb"
	"github.com/brandoyts/go-token-auth/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func loadDependencies() *appDependency {
	db, err := mongodb.NewMongodb(os.Getenv("MONGO_DATABASE_NAME"), os.Getenv("MONGO_URI"), options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})
	if err != nil {
		log.Fatal("❌ can't connect to mongodb", err)
	}

	fmt.Println("✅ successfully connected to mongodb")

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})
	pingErr := redisClient.Conn().Ping(context.Background()).Err()
	if pingErr != nil {
		log.Fatal("❌ can't connect to redis", err)
	}

	fmt.Println("✅ successfully connected to redis")

	userRepository := mongodb.NewUserRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	fmt.Println("✅ dependencies are loaded successfully")

	return &appDependency{
		db:    db,
		redis: redisClient,
		handler: &handler{
			userHandler: userHandler,
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
	defer deps.redis.Close()

	app := fiber.New()
	app.Use(logger.New(), recover.New())

	apiRouter := app.Group("/api/v1")

	apiRouter.Get("/health-check", func(c *fiber.Ctx) error {
		return c.SendString("healthy")
	})

	app.Listen(":6000")
}
