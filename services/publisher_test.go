package services

import (
	"context"
	"encoding/json"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"sync"
	"testing"
)

// MockRedis is a mock implementation of redis.Client
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	args := m.Called(ctx, channel, message)
	return args.Get(0).(*redis.IntCmd)
}

func TestPublisherService_Activate(t *testing.T) {
	service := &PublisherService{
		Lock:   &sync.Mutex{},
		Active: false,
	}

	service.Activate()
	assert.True(t, service.Active)
}

func TestPublisherService_Deactivate(t *testing.T) {
	service := &PublisherService{
		Lock:   &sync.Mutex{},
		Active: true,
	}

	service.Deactivate()
	assert.False(t, service.Active)
}

func TestPublisherService_GetStatus(t *testing.T) {
	service := &PublisherService{
		Lock:   &sync.Mutex{},
		Active: true,
	}

	status := service.GetStatus()
	assert.True(t, status)
}

func TestPublisherService_run(t *testing.T) {
	tests := []struct {
		name          string
		isActive      bool
		messages      []*repositories.Message
		publishError  error
		expectedCalls int
		shouldPublish bool
	}{
		{
			name:          "inactive service should not process messages",
			isActive:      false,
			messages:      []*repositories.Message{},
			expectedCalls: 0,
			shouldPublish: false,
		},
		{
			name:     "should process and publish messages",
			isActive: true,
			messages: []*repositories.Message{
				{Model: gorm.Model{ID: 1}, Status: repositories.MessageStatusScheduled},
				{Model: gorm.Model{ID: 2}, Status: repositories.MessageStatusScheduled},
			},
			expectedCalls: 2,
			shouldPublish: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			messageService := &MessageService{DB: db.DB}
			mockRedis := new(MockRedis)

			service := &PublisherService{
				MessageService: messageService,
				Redis:          mockRedis,
				Lock:           &sync.Mutex{},
				Context:        context.Background(),
				Active:         tt.isActive,
			}

			if tt.isActive && tt.messages != nil {
				// Pre-populate the database with test messages
				for _, msg := range tt.messages {
					err := db.DB.Create(msg).Error
					assert.NoError(t, err)
				}

				if tt.shouldPublish {
					for _, msg := range tt.messages {
						payload, _ := json.Marshal(msg)
						cmd := redis.NewIntCmd(context.Background())
						mockRedis.On("Publish", mock.Anything, "messages", payload).Return(cmd)
					}
				}
			}

			service.run(context.Background())
			mockRedis.AssertExpectations(t)
		})
	}
}
