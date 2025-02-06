package services

import (
	"context"
	"encoding/json"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
}

type PublisherService struct {
	MessageService *MessageService
	Duration       time.Duration
	Redis          Publisher
	Timer          *time.Timer
	Lock           sync.Locker
	Context        context.Context
	Active         bool
}

func (s *PublisherService) Watch() {
	s.Timer = time.NewTimer(s.Duration)
	s.run(s.Context)
	for {
		select {
		case <-s.Context.Done(): // application closed
			zap.L().Info("Publisher service stopped")
		case <-s.Timer.C: // timer expired
			s.run(s.Context)
			s.Timer.Reset(s.Duration)
		}
	}
}

func (s *PublisherService) Activate() {
	zap.L().Debug("Activating publisher service")
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Active = true
}

func (s *PublisherService) Deactivate() {
	zap.L().Debug("Deactivating publisher service")
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Active = false
}

func (s *PublisherService) GetStatus() bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.Active
}

func (s *PublisherService) run(ctx context.Context) {
	zap.L().Debug("Checking messages to process")
	if !s.Active {
		zap.L().Debug("Publisher service is not active")
		return
	}
	result, err := s.MessageService.List(ctx, 2, 0, "id desc", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", repositories.MessageStatusScheduled)
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
