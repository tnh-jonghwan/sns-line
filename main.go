package main

import (
	"messaging-line/app"
	"messaging-line/config"
	"messaging-line/domain/webhook"
	"messaging-line/jwt"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			app.NewFiberApp,
			config.GetEnv,
			jwt.GetAccessToken,
		),

		webhook.WebhookModule,
	).Run()
}
