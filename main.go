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
		// Config
		fx.Provide(config.GetEnv),

		// JWT
		jwt.Module,

		// Webhook domain
		webhook.Module,

		// Fiber app
		fx.Provide(app.NewFiberApp),

		// Routes
		fx.Invoke(webhook.RegisterRoutes),
	).Run()
}
