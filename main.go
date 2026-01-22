package main

import (
	"log"
	"os"
	"sns/app"
	"sns/config"
	"sns/domain/eventHub"
	lineClient "sns/domain/line"
	instagramWebhook "sns/domain/webhook/instagram"
	lineWebhook "sns/domain/webhook/line"

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
