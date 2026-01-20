package app

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// NewFiberApp - Fiber 앱 생성 및 라이프사이클 관리
func NewFiberApp(lc fx.Lifecycle) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

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
