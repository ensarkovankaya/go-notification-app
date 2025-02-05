package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (h *AppHandler) Setup(router fiber.Router) {
	router.Get("/readyz", h.Ready)
	router.Get("/healthz", h.Health)
}

// Ready is a handler function that returns OK if the application is ready
func (h *AppHandler) Ready(c *fiber.Ctx) error {
	if err := h.checkDatabaseConnection(c.UserContext()); err != nil {
		return err
	}
	if err := h.checkRedisConnection(c.UserContext()); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"status": "OK"})
}

// Health is a handler function that returns OK if the application is healthy
func (h *AppHandler) Health(c *fiber.Ctx) error {
	if err := h.checkDatabaseConnection(c.UserContext()); err != nil {
		return err
	}
	if err := h.checkRedisConnection(c.UserContext()); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"status": "OK"})
}

// checkDatabaseConnection is a helper function that check the database connection is valid
func (h *AppHandler) checkDatabaseConnection(ctx context.Context) error {
	db, err := h.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

func (h *AppHandler) checkRedisConnection(ctx context.Context) error {
	if err := h.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}
	return nil
}
