package services

import (
	"context"
	"encoding/json"
	"github.com/ensarkovankaya/go-notification-app/clients"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SubscriberService struct {
	MessageService *MessageService
	Redis          *redis.Client
	WebhookClient  *clients.WebhookClient
	Context        context.Context
}

func (s *SubscriberService) Watch() {
	subscriber := s.Redis.Subscribe(s.Context, "messages")
	for {
		select {
		case <-s.Context.Done():
			return
		default:
			msg, err := subscriber.ReceiveMessage(s.Context)
			if err != nil {
				zap.L().Error("Failed to receive message", zap.Error(err))
				continue
			}
			s.run(s.Context, msg.Payload)
		}
	}
}

func (s *SubscriberService) run(ctx context.Context, payload string) {
	zap.L().Debug("Message received", zap.String("payload", payload))
	// Unmarshall message
	message := new(repositories.Message)
	err := json.Unmarshal([]byte(payload), message)
	if err != nil {
		zap.L().Error("Failed to Unmarshall payload", zap.Error(err))
		return
	}

	// Send to client
	messageID, err := s.WebhookClient.Send(ctx, clients.Payload{
		To:      message.Recipient,
		Content: message.Content,
	})
	if err != nil {
		_ = s.MessageService.MarkAsFailed(ctx, message.ID)
	} else {
		_ = s.MessageService.MarkAsSend(ctx, message.ID, messageID)
	}
}
