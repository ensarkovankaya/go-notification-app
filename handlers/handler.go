package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AppHandler struct {
	DB *gorm.DB
}

func (h *AppHandler) Setup(router fiber.Router) {
	router.Get("/readyz", h.Ready)
	router.Get("/healthz", h.Health)
}

// Ready is a handler function that returns OK if the application is ready
func (h *AppHandler) Ready(c *fiber.Ctx) error {
	if err := h.validateDatabase(c.UserContext()); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).Send([]byte("OK"))
}

// Health is a handler function that returns OK if the application is healthy
func (h *AppHandler) Health(c *fiber.Ctx) error {
	if err := h.validateDatabase(c.UserContext()); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).Send([]byte("OK"))
}

func (h *AppHandler) validateDatabase(ctx context.Context) error {
	db, err := h.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}
