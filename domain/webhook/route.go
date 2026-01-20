package webhook

import (
	"messaging-line/jwt"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes - Webhook 라우팅 설정
func RegisterRoutes(app *fiber.App, handler *WebhookHandler) {
	baseRouter := app.Group("/")

	baseRouter.Get("/", healthCheck)
	baseRouter.Post("/webhook", handler.Handle)
	baseRouter.Get("/refresh-token", refreshTokenHandler)
}

// healthCheck - 헬스체크 핸들러
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "LINE Webhook Server",
		"status":  "running",
	})
}

// refreshTokenHandler - Access Token 갱신 핸들러
func refreshTokenHandler(c *fiber.Ctx) error {
	newToken := jwt.GetAccessToken()
	return c.JSON(fiber.Map{
		"message":      "Token refreshed",
		"access_token": newToken,
	})
}
