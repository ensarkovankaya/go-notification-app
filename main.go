package main

import (
	"fmt"
	"github.com/ensarkovankaya/go-messagingapp/pkg"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() { _ = pkg.Logger.Sync() }()
	// Initialize Fiber
	app := fiber.New(fiber.Config{
		IdleTimeout:  pkg.Cnf.IdleTimeout,
		ReadTimeout:  pkg.Cnf.ReadTimeout,
		WriteTimeout: pkg.Cnf.WriteTimeout,
	})

	// Run http server
	go func() {
		address := fmt.Sprintf(":%s", pkg.Cnf.Port)
		pkg.Logger.Info(fmt.Sprintf("Application listening at %v", address))
		if err := app.Listen(address); err != nil {
			pkg.Logger.Error("Application could not started", zap.Error(err))
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Close application
	if err := app.Shutdown(); err != nil {
		pkg.Logger.Error("Server shutdown failed", zap.Error(err))
	} else {
		pkg.Logger.Info("Server shutdown succeeded")
	}
}
