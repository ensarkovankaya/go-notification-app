package repositories

import "fmt"

type MessageStatus string

var (
	StatusScheduled MessageStatus = "SCHEDULED"
	StatusSuccess   MessageStatus = "SUCCESS"
	StatusFailed    MessageStatus = "FAILED"
)

// Scan implements the sql.Scanner interface for enum values in gorm
func (s *MessageStatus) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid status: %v", value)
	}
	*s = MessageStatus(v)
	return nil
}

// Value implements the driver.Valuer interface for enum values in gorm
func (s *MessageStatus) Value() (interface{}, error) {
	return string(*s), nil
}
