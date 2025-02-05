package handlers

import (
	"github.com/ensarkovankaya/go-notification-app/common"
	"github.com/ensarkovankaya/go-notification-app/models"
	"github.com/ensarkovankaya/go-notification-app/repositories"
	"github.com/ensarkovankaya/go-notification-app/services"
	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MessageHandler struct {
	MessageService *services.MessageService
}

func (h *MessageHandler) Setup(router fiber.Router) {
	router.Get("/messages", h.ListMessages)
	router.Post("/messages", h.SendMessage)
}

func (h *MessageHandler) ListMessages(c *fiber.Ctx) error {
	limit := int64(c.QueryInt("limit", 100))
	offset := int64(c.QueryInt("offset", 0))
	order := c.Query("order", "id desc")

	var filters []func(db *gorm.DB) *gorm.DB
	if c.Query("status") != "" {
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", c.Query("status"))
		})
	}
	if c.Query("recipient") != "" {
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("recipient = ?", c.Query("recipient"))
		})
	}
	result, err := h.MessageService.List(c.UserContext(), limit, offset, order, filters...)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIError{Code: 1000, Description: err.Error()})
	}
	// Convert result to api response
	response := &models.MessageList{
		PaginatedResponse: models.PaginatedResponse{
			Limit:  &limit,
			Offset: &offset,
			Total:  &result.Total,
		},
	}
	for _, message := range result.Data {
		response.Data = append(response.Data, h.serializeMessage(message))
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
	request := new(models.CreateMessageRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIError{Code: 1000, Description: err.Error()})
	}
	if err := request.Validate(strfmt.Default); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIError{Code: 1000, Description: err.Error()})
	}
	message, err := h.MessageService.Create(c.UserContext(), *request.Recipient, *request.Content)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIError{Code: 1001, Description: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(h.serializeMessage(message))
}

func (h *MessageHandler) serializeMessage(message *repositories.Message) *models.Message {
	return &models.Message{
		ID:        int64(message.ID),
		Content:   message.Content,
		Recipient: message.Recipient,
		SentTime:  common.SqlNullTimeToPtr(message.SendTime),
		MessageID: common.SqlNullStringToPtr(message.MessageID),
		Status:    string(message.Status),
		CreatedAt: strfmt.DateTime(message.CreatedAt),
		UpdatedAt: strfmt.DateTime(message.UpdatedAt),
	}
}
