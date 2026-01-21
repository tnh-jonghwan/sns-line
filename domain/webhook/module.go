package webhook

import "go.uber.org/fx"

var WebhookModule = fx.Options(
	fx.Provide(
		NewLineClient,
		NewWebhookService,
		NewWebhookHandler,
	),
	fx.Invoke(
		RegisterWebhookRoutes,
	),
)
