package services

import (
	"context"
	"fmt"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

// Create creates a new message with given recipient and content
func (s *MessageService) Create(ctx context.Context, recipient string, content string) (*repositories.Message, error) {
	zap.L().Info("Creating message", zap.String("recipient", recipient), zap.String("content", content))
	message := &repositories.Message{
		Recipient: recipient,
		Content:   content,
		Status:    repositories.StatusScheduled,
	}

	if err := s.DB.WithContext(ctx).Create(message).Error; err != nil {
		zap.L().Error("Failed to create message", zap.Error(err))
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	zap.L().Debug("Message created", zap.Uint("id", message.ID))
	return message, nil
}

// List fetch messages with given filters from database
func (s *MessageService) List(ctx context.Context, limit, offset int64, order string, filters ...func(db *gorm.DB) *gorm.DB) (*repositories.MessageList, error) {
	zap.L().Info("Listing messages", zap.Int64("limit", limit), zap.Int64("offset", offset), zap.String("order", order))
	var total int64
	var messages []*repositories.Message
	db := s.DB.WithContext(ctx)
	for _, query := range filters {
		db = query(db)
	}
	if err := db.Model(&repositories.Message{}).Count(&total).Error; err != nil {
		zap.L().Error("Failed to count messages", zap.Error(err))
		return nil, fmt.Errorf("failed to count messages: %w", err)
	}

	if err := db.Limit(int(limit)).Offset(int(offset)).Order(order).Find(&messages).Error; err != nil {
		zap.L().Error("Failed to fetch messages", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	return &repositories.MessageList{
		Limit:  limit,
		Offset: offset,
		Total:  total,
		Data:   messages,
	}, nil
}

func (s *MessageService) MarkAsSend(ctx context.Context, id uint, messageID string) error {
	zap.L().Info("Marking message as sent", zap.Uint("id", id), zap.String("message_id", messageID))
	if err := s.DB.WithContext(ctx).Model(&repositories.Message{}).Where("id = ?", id).UpdateColumns(map[string]any{
		"status":     repositories.StatusSent,
		"message_id": messageID,
		"send_time":  s.DB.NowFunc(),
		"updated_at": s.DB.NowFunc(),
	}).Error; err != nil {
		zap.L().Error("Failed to mark message as sent", zap.Error(err))
		return fmt.Errorf("failed to mark message as sent: %w", err)
	}
	zap.L().Debug("Message marked as sent", zap.Uint("id", id))
	return nil
}

func (s *MessageService) MarkAsFailed(ctx context.Context, id uint) error {
	zap.L().Info("Marking message as failed", zap.Uint("id", id))
	if err := s.DB.WithContext(ctx).Model(&repositories.Message{}).Where("id = ?", id).UpdateColumns(map[string]any{
		"status":     repositories.StatusFailed,
		"updated_at": s.DB.NowFunc(),
	}).Error; err != nil {
		zap.L().Error("Failed to mark message as failed", zap.Error(err))
		return fmt.Errorf("failed to mark message as failed: %w", err)
	}
	zap.L().Debug("Message marked as failed", zap.Uint("id", id))
	return nil
}
