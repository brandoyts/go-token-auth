package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/brandoyts/go-token-auth/internal/user"
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

	requestBody := new(LoginRequest)
	err := c.BodyParser(requestBody)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	authToken, err := h.authService.Login(ctx, user.User{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	})
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	return c.JSON(LoginResponse{
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
	})

}
