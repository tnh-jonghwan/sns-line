package main

import (
	"messaging-line/app"
	"messaging-line/config"
	"messaging-line/domain"
	"messaging-line/domain/webhook"
	"messaging-line/jwt"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.GetEnv),
		fx.Provide(jwt.GetAccessToken),

		// Domain modules
		webhook.Module,       // WebhookHandler 제공
		domain.HandlerModule, // []Handler 제공

		// App initialization
		fx.Invoke(app.NewApp),
	).Run()
}
