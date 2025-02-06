package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ensarkovankaya/go-notification-app/clients"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

// MockSubscriber mocks the Subscriber interface
type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	args := m.Called(ctx, channels)
	return args.Get(0).(*redis.PubSub)
}

// MockSMSClient mocks the SMSClient interface
type MockSMSClient struct {
	mock.Mock
}

func (m *MockSMSClient) Send(ctx context.Context, payload clients.Payload) (string, error) {
	args := m.Called(ctx, payload)
	return args.String(0), args.Error(1)
}

func TestSubscriberService_run(t *testing.T) {
	tests := []struct {
		name           string
		message        *repositories.Message
		smsResponse    string
		smsError       error
		expectedStatus repositories.MessageStatus
	}{
		{
			name: "successful message delivery",
			message: &repositories.Message{
				Model:     gorm.Model{ID: 1},
				Content:   "Test message",
				Recipient: "+1234567890",
				Status:    repositories.MessageStatusScheduled,
			},
			smsResponse:    "msg_123",
			smsError:       nil,
			expectedStatus: repositories.MessageStatusSuccess,
		},
		{
			name: "failed message delivery",
			message: &repositories.Message{
				Model:     gorm.Model{ID: 2},
				Content:   "Test message",
				Recipient: "+1234567890",
				Status:    repositories.MessageStatusScheduled,
			},
			smsResponse:    "",
			smsError:       errors.New("sms delivery failed"),
			expectedStatus: repositories.MessageStatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db := setupTestDB(t)
			messageService := &MessageService{DB: db.DB}
			mockSMSClient := new(MockSMSClient)
			mockSubscriber := new(MockSubscriber)

			service := &SubscriberService{
				MessageService: messageService,
				Redis:          mockSubscriber,
				SmsClient:      mockSMSClient,
				Context:        context.Background(),
			}

			// Create test message in DB
			err := db.DB.Create(tt.message).Error
			assert.NoError(t, err)

			// Setup mock expectations
			mockSMSClient.On("Send", mock.Anything, clients.Payload{
				To:      tt.message.Recipient,
				Content: tt.message.Content,
			}).Return(tt.smsResponse, tt.smsError)

			// Convert message to JSON payload
			payload, err := json.Marshal(tt.message)
			assert.NoError(t, err)

			// Run the test
			service.run(context.Background(), string(payload))

			// Verify the message status was updated correctly
			var updatedMessage repositories.Message
			err = db.DB.First(&updatedMessage, tt.message.ID).Error
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, updatedMessage.Status)

			if tt.smsError == nil {
				assert.Equal(t, tt.smsResponse, updatedMessage.MessageID.String)
				assert.True(t, updatedMessage.SendTime.Valid)
			}

			mockSMSClient.AssertExpectations(t)
		})
	}
}

func TestSubscriberService_run_InvalidPayload(t *testing.T) {
	// Setup
	messageService := &MessageService{DB: setupTestDB(t).DB}
	mockSMSClient := new(MockSMSClient)
	mockSubscriber := new(MockSubscriber)

	service := &SubscriberService{
		MessageService: messageService,
		Redis:          mockSubscriber,
		SmsClient:      mockSMSClient,
		Context:        context.Background(),
	}

	// Run with invalid JSON payload
	service.run(context.Background(), "invalid json payload")

	// Verify that SMS client was never called
	mockSMSClient.AssertNotCalled(t, "Send")
}
