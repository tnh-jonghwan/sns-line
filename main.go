package main

import (
	"log"
	"os"
	"sns-line/app"
	"sns-line/config"
	"sns-line/domain/webhook"
	"sns-line/jwt"

	"go.uber.org/fx"
)

func main() {
	// log를 stdout으로 출력
	log.SetOutput(os.Stdout)

	fx.New(
		fx.Provide(
			app.NewFiberApp,
			config.GetEnv,
			jwt.GetAccessToken,
		),

		webhook.WebhookModule,
	).Run()
}
