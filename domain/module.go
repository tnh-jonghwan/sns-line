package domain

import (
	"messaging-line/domain/webhook"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// Handler - 모든 domain handler가 구현해야 하는 인터페이스
type Handler interface {
	RegisterRoutes(baseRouter fiber.Router)
}

// GetAllHandlers - 모든 handler를 통합 (fx가 주입)
func GetAllHandlers(
	webhook *webhook.WebhookHandler,
	// 나중에 추가:
	// auth *auth.AuthHandler,
	// user *user.UserHandler,
) []Handler {
	return []Handler{
		webhook,
	}
}

// HandlerModule - Handler 통합 모듈
var HandlerModule = fx.Module("handler",
	fx.Provide(GetAllHandlers),
)
