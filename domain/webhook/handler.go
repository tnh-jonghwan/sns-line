package webhook

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type WebhookHandler struct {
	service *WebhookService
}

func NewWebhookHandler(service *WebhookService) *WebhookHandler {
	return &WebhookHandler{
		service: service,
	}
}

func (h *WebhookHandler) Handle(c *fiber.Ctx) error {
	// 공통: Request Body 파싱
	var req WebhookRequest

	log.Println("webhook handle 호출")

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Service layer에서 이벤트 처리
	if err := h.service.HandleEvents(req.Events); err != nil {
		log.Printf("Error handling events: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process events",
		})
	}

	// LINE은 반드시 200 OK를 받아야 함
	return c.SendStatus(fiber.StatusOK)
}
