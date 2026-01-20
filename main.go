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
		fx.Provide(
			config.GetEnv,
			jwt.GetAccessToken,

			webhook.NewWebhookHandler,
		),

		// Domain handler
		domain.HandlerModule,

		// App initialization
		fx.Invoke(app.NewApp),
	).Run()
}
