package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/brandoyts/go-token-auth/internal/infrastructure/jwtAuth"
	"github.com/brandoyts/go-token-auth/internal/infrastructure/redisClient"
	"github.com/gofiber/fiber/v2"
)

func authChecker(redisClient *redisClient.RedisClient, jwtAuth *jwtAuth.JwtAuth) fiber.Handler {
	return func(c *fiber.Ctx) error {

		header := c.Get("Authorization")

		authHeader := strings.Split(header, " ")

		if len(authHeader) != 2 {
			return c.SendStatus(http.StatusUnauthorized)
		}

		if strings.ToLower(authHeader[0]) != "bearer" {
			return c.SendStatus(http.StatusUnauthorized)
		}

		accessToken := authHeader[1]

		err := jwtAuth.Verify(accessToken)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(http.StatusUnauthorized)
		}

		// check if token is blacklisted
		cache, err := redisClient.Get(accessToken)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(http.StatusUnauthorized)
		}

		if cache != "" {
			fmt.Println("trying to use a token that has been blacklisted")
			return c.SendStatus(http.StatusUnauthorized)
		}

		c.Locals("access_token", accessToken)

		return c.Next()
	}
}
