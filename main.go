package main

import (
	"fmt"
	"github.com/ensarkovankaya/go-message-broker/common"
	"github.com/ensarkovankaya/go-message-broker/handlers"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Handle application panic
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("Application panicked", zap.Any("error", err))
			_ = zap.L().Sync()
		}
	}()
	defer func() { _ = zap.L().Sync() }()

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		IdleTimeout:  common.Cnf.IdleTimeout,
		ReadTimeout:  common.Cnf.ReadTimeout,
		WriteTimeout: common.Cnf.WriteTimeout,
	})

	// Handlers
	rootRouter := app.Group("/api")
	appHandler := handlers.AppHandler{DB: DB}
	appHandler.Setup(rootRouter)

	// Run http server
	go func() {
		address := fmt.Sprintf(":%s", common.Cnf.Port)
		zap.L().Info(fmt.Sprintf("Application listening at %v", address))
		if err := app.Listen(address); err != nil {
			zap.L().Error("Application could not started", zap.Error(err))
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Close http application
	if err := app.Shutdown(); err != nil {
		zap.L().Error("Server shutdown failed", zap.Error(err))
	} else {
		zap.L().Info("Server shutdown succeeded")
	}
}
