package webhook

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterWebhookRoutes(
	app *fiber.App,
	webhookHandler *WebhookHandler,
) {
	app.Post("/webhook", webhookHandler.Handle)
}
