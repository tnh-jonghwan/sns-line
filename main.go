package main

import (
	"log"
	"os"
	"sns-line/app"
	"sns-line/config"
	"sns-line/domain/eventHub"
	lineClient "sns-line/domain/line"
	instagramWebhook "sns-line/domain/webhook/instagram"
	lineWebhook "sns-line/domain/webhook/line"

	"go.uber.org/fx"
)

func main() {
	// log를 stdout으로 출력
	log.SetOutput(os.Stdout)

	fx.New(
		fx.Provide(
			app.NewFiberApp,
			config.GetEnv,
		),

		lineClient.LineModule,
		eventHub.EventHubModule,
		lineWebhook.LineWebhookModule,
		instagramWebhook.InstagramModule,
	).Run()
}
