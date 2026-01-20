package main

import (
	"fmt"
	"log"
	"messaging-line/config"
	"messaging-line/domain/webhook"
	"messaging-line/jwt"

	"github.com/gofiber/fiber/v2"
)

// 환경 변수 로드
var env = config.GetEnv()

func main() {
	// Access Token 생성
	accessToken := jwt.GetAccessToken()
	fmt.Println("Access Token:", accessToken)

	// 의존성 주입 (Dependency Injection)
	lineClient := webhook.NewLineClient(accessToken, env.LineApiPrefix)
	webhookService := webhook.NewWebhookService(lineClient)
	webhookHandler := webhook.NewWebhookHandler(webhookService)

	// Fiber 앱 생성
	app := fiber.New()

	// 라우팅 설정
	BuildRoute(app, webhookHandler)

	// 서버 시작
	log.Fatal(app.Listen(":3000"))
}

func BuildRoute(app *fiber.App, webhookHandler *webhook.WebhookHandler) {
	baseRouter := app.Group("/")

	baseRouter.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "LINE Webhook Server",
			"status":  "running",
		})
	})

	baseRouter.Post("/webhook", webhookHandler.Handle)

	baseRouter.Get("/refresh-token", func(c *fiber.Ctx) error {
		newToken := jwt.GetAccessToken()
		return c.JSON(fiber.Map{
			"message":      "Token refreshed",
			"access_token": newToken,
		})
	})
}
