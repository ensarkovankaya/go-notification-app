package repositories

import (
	"database/sql"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Recipient string         `gorm:"size:256;not null;" json:"recipient,omitempty"`
	Content   string         `gorm:"size:2048;not null;" json:"content,omitempty"`
	Status    MessageStatus  `gorm:"size:50;not null;index;" json:"status,omitempty"`
	MessageID sql.NullString `json:"message_id"`
	SendTime  sql.NullTime   `json:"send_time"`
}

type MessageList struct {
	Limit  int64
	Offset int64
	Total  int64
	Data   []*Message
}
