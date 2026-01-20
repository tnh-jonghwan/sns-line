package app

import (
	"context"
	"log"

	"messaging-line/domain"
	"messaging-line/jwt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// NewApp - Fiber 앱 생성 및 라우팅 설정
func NewApp(lc fx.Lifecycle, handlers []domain.Handler) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

	// Base router 생성
	baseRouter := app.Group("/")

	// 전역 라우트 등록
	baseRouter.Get("/refresh-token", refreshTokenHandler)

	// 각 domain handler의 라우트 등록
	for _, h := range handlers {
		h.RegisterRoutes(baseRouter)
	}

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
}

// refreshTokenHandler - Access Token 갱신 핸들러
func refreshTokenHandler(c *fiber.Ctx) error {
	newToken := jwt.GetAccessToken()
	return c.JSON(fiber.Map{
		"message":      "Token refreshed",
		"access_token": newToken,
	})
}
