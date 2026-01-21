package app

import (
	"context"
	"log"

	"messaging-line/jwt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// NewFiberApp - Fiber 앱 생성, 전역 라우트 등록, 서버 시작
func NewFiberApp(lc fx.Lifecycle) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

	// 전역 라우트 등록
	buildGlobalRoutes(app)

	// 서버 라이프사이클 관리
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(":3000"); err != nil {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			log.Println("Server started on :3000")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return app.Shutdown()
		},
	})

	return app
}

func buildGlobalRoutes(app *fiber.App) {
	// test
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, SNS-LINE SERVER!")
	})

	// refresh
	app.Get("/refresh-token", func(c *fiber.Ctx) error {
		newToken := jwt.GetAccessToken()
		return c.JSON(fiber.Map{
			"message":      "Token refreshed",
			"access_token": newToken,
		})
	})
}
