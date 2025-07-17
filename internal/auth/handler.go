package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/brandoyts/go-token-auth/internal/shared"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ipAddress := c.IP()

	requestBody := new(LoginRequest)
	err := c.BodyParser(requestBody)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	authToken, err := h.authService.Login(ctx, LoginInput{
		IPAddress: ipAddress,
		Email:     requestBody.Email,
		Password:  requestBody.Password,
	})
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	c.Cookie(&fiber.Cookie{
		Name:     shared.COOKIES_REFRESH_TOKEN,
		Value:    authToken.RefreshToken,
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(TokenResponse{
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
	})
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	refreshToken := c.Cookies(shared.COOKIES_REFRESH_TOKEN)
	if refreshToken == "" {
		return c.SendStatus(http.StatusUnauthorized)
	}

	ipAddress := c.IP()

	// refactor this line later
	accessToken := strings.Split(c.Get("Authorization"), " ")[1]

	result, err := h.authService.RefreshToken(ctx, RefreshTokenInput{
		IPAddress:    ipAddress,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(http.StatusUnauthorized)
	}

	c.Cookie(&fiber.Cookie{
		Name:     shared.COOKIES_REFRESH_TOKEN,
		Value:    result.RefreshToken,
		HTTPOnly: true,
	})

	return c.JSON(TokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	accessToken, ok := c.Locals("access_token").(string)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}

	refreshToken := c.Cookies(shared.COOKIES_REFRESH_TOKEN)
	if refreshToken == "" {
		return c.SendStatus(http.StatusUnauthorized)
	}

	err := h.authService.Logout(ctx, accessToken, refreshToken)
	if err != nil {
		log.Fatal(err)
		return c.SendStatus(http.StatusBadGateway)
	}

	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:     shared.COOKIES_REFRESH_TOKEN,
		Value:    "",
		Expires:  expired,
		HTTPOnly: true,
	})

	return c.SendStatus(http.StatusNoContent)
}
