package handlers

import (
	"github.com/ensarkovankaya/go-notification-app/models"
	"github.com/ensarkovankaya/go-notification-app/services"
	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber/v2"
)

type CronHandler struct {
	PublisherService *services.PublisherService
}

func (h *CronHandler) Setup(router fiber.Router) {
	router.Get("/cron", h.GetStatus)
	router.Post("/cron", h.UpdateStatus)
}

func (h *CronHandler) GetStatus(c *fiber.Ctx) error {
	status := h.PublisherService.GetStatus()
	return c.Status(fiber.StatusOK).JSON(models.CronStatus{Active: &status})
}

func (h *CronHandler) UpdateStatus(c *fiber.Ctx) error {
	request := new(models.CronStatus)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIError{Code: 1000, Description: err.Error()})
	}
	if err := request.Validate(strfmt.Default); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.APIError{Code: 1000, Description: err.Error()})
	}
	if *request.Active {
		h.PublisherService.Activate()
	} else {
		h.PublisherService.Deactivate()
	}
	return c.Status(fiber.StatusOK).JSON(models.CronStatus{Active: request.Active})
}
