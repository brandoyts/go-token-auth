package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	userService *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) FindUserById(c *fiber.Ctx) error {
	userId := c.Params("id")

	if userId == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := h.userService.FindUserById(ctx, userId)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (h *Handler) FindUser(c *fiber.Ctx) error {
	requestBody := new(User)

	err := c.BodyParser(requestBody)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := h.userService.FindUser(ctx, *requestBody)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.JSON(result)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	requestBody := new(User)

	err := c.BodyParser(requestBody)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := h.userService.CreateUser(ctx, *requestBody)
	if err != nil {
		return err
	}

	return c.SendString(result)
}
