package main

import (
	"fmt"
	"github.com/ensarkovankaya/go-notification-app/handlers"
	"github.com/ensarkovankaya/go-notification-app/services"
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

	// Initialize HTTP Server
	app := fiber.New(fiber.Config{
		IdleTimeout:  Cnf.IdleTimeout,
		ReadTimeout:  Cnf.ReadTimeout,
		WriteTimeout: Cnf.WriteTimeout,
	})

	// Services
	messageService := &services.MessageService{DB: DB}

	// Handlers
	rootRouter := app.Group("/api")
	appHandler := handlers.AppHandler{DB: DB}
	appHandler.Setup(rootRouter)

	messageHandler := handlers.MessageHandler{MessageService: messageService}
	messageHandler.Setup(rootRouter)

	// Start http server
	go func() {
		address := fmt.Sprintf(":%s", Cnf.Port)
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
