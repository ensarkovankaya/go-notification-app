package repositories

import (
	"database/sql"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Recipient string        `gorm:"size:256;not null;"`
	Content   string        `gorm:"size:2048;not null;"`
	Status    MessageStatus `gorm:"size:50;not null;index;"`
	MessageID sql.NullString
	SendTime  sql.NullTime
}
