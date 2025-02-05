package services

import (
	"context"
	"encoding/json"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type PublisherService struct {
	MessageService *MessageService
	Duration       time.Duration
	Redis          *redis.Client
}

func (s *PublisherService) Start(ctx context.Context) {
	s.run(ctx)
	timer := time.NewTimer(s.Duration)
	for {
		select {
		case <-ctx.Done():
			zap.L().Info("Publisher service stopped")
			timer.Stop()
			return
		case <-timer.C:
			s.run(ctx)
			timer.Reset(s.Duration)
		}
	}
}

func (s *PublisherService) run(ctx context.Context) {
	zap.L().Debug("Checking messages to process")
	result, err := s.MessageService.List(ctx, 2, 0, "id desc", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", repositories.StatusScheduled)
	})
	if err != nil {
		zap.L().Error("Failed to fetch messages", zap.Error(err))
		return
	}
	zap.L().Debug("Message found", zap.Any("count", len(result.Data)))

	var payload []byte
	for _, message := range result.Data {
		if payload, err = json.Marshal(message); err != nil {
			zap.L().Error("Failed to marshall message", zap.Error(err))
			continue
		}
		if err = s.Redis.Publish(ctx, "messages", payload).Err(); err != nil {
			zap.L().Error("Failed to publish message", zap.Any("message", message))
		} else {
			zap.L().Debug("Message sent", zap.Any("message", message))
		}
	}
}
