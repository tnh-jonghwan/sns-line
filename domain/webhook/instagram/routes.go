package instagram

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterInstagramRoutes(app *fiber.App, handler *InstagramHandler) {
	// GET /webhook/instagram - 웹훅 검증
	app.Get("/webhook/instagram", handler.HandleVerify)

	// POST /webhook/instagram - 웹훅 수신
	app.Post("/webhook/instagram", handler.HandleWebhook)
}
