package main

import (
	"context"
	"fmt"
	"github.com/ensarkovankaya/go-notification-app/clients"
	"github.com/ensarkovankaya/go-notification-app/handlers"
	"github.com/ensarkovankaya/go-notification-app/services"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
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
	defer closeDB()
	defer closeRedis()
	defer func() { _ = zap.L().Sync() }()
	ctx := context.Background() // application context

	// Initialize HTTP Server
	app := fiber.New(fiber.Config{
		IdleTimeout:  Cnf.IdleTimeout,
		ReadTimeout:  Cnf.ReadTimeout,
		WriteTimeout: Cnf.WriteTimeout,
	})

	// Clients
	webhookClient := clients.NewWebhookClient(Cnf.WebhookID)

	// Services
	messageService := &services.MessageService{DB: DB}
	publisherService := &services.PublisherService{
		MessageService: messageService,
		Duration:       Cnf.CronTTL,
		Redis:          Redis,
		Context:        ctx,
		Lock:           &sync.Mutex{},
		Active:         true,
	}
	subscriberService := &services.SubscriberService{
		MessageService: messageService,
		Redis:          Redis,
		WebhookClient:  webhookClient,
		Context:        ctx,
	}

	// Handlers
	rootRouter := app.Group("/api")
	appHandler := handlers.AppHandler{DB: DB, Redis: Redis}
	appHandler.Setup(rootRouter)

	messageHandler := handlers.MessageHandler{MessageService: messageService}
	messageHandler.Setup(rootRouter)

	cronHandler := handlers.CronHandler{PublisherService: publisherService}
	cronHandler.Setup(rootRouter)

	// Start http server
	go func() {
		address := fmt.Sprintf(":%s", Cnf.Port)
		zap.L().Info(fmt.Sprintf("Application listening at %v", address))
		if err := app.Listen(address); err != nil {
			zap.L().Error("Application could not started", zap.Error(err))
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	// Start Publisher cron job
	go func() {
		defer wg.Done()
		publisherService.Watch()
	}()
	// Start Subscriber
	go func() {
		defer wg.Done()
		subscriberService.Watch()
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // listens for interrupt signal

	// Stop accepting requests and wait for all requests to finish
	if err := app.Shutdown(); err != nil {
		zap.L().Error("Server shutdown failed", zap.Error(err))
	} else {
		zap.L().Info("Server shutdown succeeded")
	}

	ctx.Done() // Tell publisher and subscriber to stop running
	wg.Wait()  // wait publisher and subscriber to stop
}
