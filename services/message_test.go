package services

import (
	"context"
	"database/sql"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type testDB struct {
	DB *gorm.DB
}

func setupTestDB(t *testing.T) *testDB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&repositories.Message{})
	assert.NoError(t, err)

	return &testDB{DB: db}
}

func (tdb *testDB) insertMessages(t *testing.T) []*repositories.Message {
	messages := []*repositories.Message{
		{
			Content:   "Test message 1",
			Recipient: "+1234567890",
			Status:    repositories.MessageStatusScheduled,
		},
		{
			Content:   "Test message 2",
			Recipient: "+9876543210",
			Status:    repositories.MessageStatusSuccess,
			SendTime:  sql.NullTime{Time: time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC), Valid: true},
			MessageID: sql.NullString{String: "message_id_1", Valid: true},
		},
		{
			Content:   "Test message 3",
			Recipient: "+1234567890",
			Status:    repositories.MessageStatusFailed,
			SendTime:  sql.NullTime{Time: time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC), Valid: true},
		},
	}

	result := tdb.DB.Create(&messages)
	assert.NoError(t, result.Error)
	return messages
}

func TestMessageService_List(t *testing.T) {
	t.Run("test list with default parameters", func(t *testing.T) {
		tdb := setupTestDB(t)
		seededMessages := tdb.insertMessages(t)

		service := MessageService{DB: tdb.DB}
		result, err := service.List(context.Background(), 100, 0, "id desc")

		assert.NoError(t, err)
		assert.Equal(t, int64(len(seededMessages)), result.Total)
		assert.Len(t, result.Data, len(seededMessages))
	})

	t.Run("test list with status filter", func(t *testing.T) {
		tdb := setupTestDB(t)
		tdb.insertMessages(t)

		service := MessageService{DB: tdb.DB}
		filter := func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", repositories.MessageStatusScheduled)
		}

		result, err := service.List(context.Background(), 100, 0, "id desc", filter)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.Total)
		assert.Len(t, result.Data, 1)
		assert.Equal(t, repositories.MessageStatusScheduled, result.Data[0].Status)
	})

	t.Run("test list with recipient filter", func(t *testing.T) {
		tdb := setupTestDB(t)
		tdb.insertMessages(t)

		service := MessageService{DB: tdb.DB}
		filter := func(db *gorm.DB) *gorm.DB {
			return db.Where("recipient = ?", "+1234567890")
		}

		result, err := service.List(context.Background(), 100, 0, "id desc", filter)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), result.Total)
		assert.Len(t, result.Data, 2)
		for _, msg := range result.Data {
			assert.Equal(t, "+1234567890", msg.Recipient)
		}
	})

	t.Run("test list with pagination", func(t *testing.T) {
		tdb := setupTestDB(t)
		tdb.insertMessages(t)

		service := MessageService{DB: tdb.DB}
		result, err := service.List(context.Background(), 1, 1, "id desc")

		assert.NoError(t, err)
		assert.Equal(t, int64(3), result.Total) // Total should still be 3
		assert.Len(t, result.Data, 1)           // But only 1 result returned
	})

	t.Run("test list with ordering", func(t *testing.T) {
		tdb := setupTestDB(t)
		tdb.insertMessages(t)

		service := MessageService{DB: tdb.DB}
		result, err := service.List(context.Background(), 100, 0, "id asc")

		assert.NoError(t, err)
		assert.Equal(t, int64(3), result.Total)
		assert.Len(t, result.Data, 3)
		// Verify ordering
		assert.Less(t, result.Data[0].ID, result.Data[1].ID)
		assert.Less(t, result.Data[1].ID, result.Data[2].ID)
	})
}

func TestMessageService_Create(t *testing.T) {
	t.Run("test create message", func(t *testing.T) {
		tdb := setupTestDB(t)
		service := MessageService{DB: tdb.DB}

		message, err := service.Create(context.Background(), "+1234567890", "Test content")

		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, "+1234567890", message.Recipient)
		assert.Equal(t, "Test content", message.Content)
		assert.Equal(t, repositories.MessageStatusScheduled, message.Status)
	})
}
