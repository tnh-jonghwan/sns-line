package webhook

import (
	"go.uber.org/fx"
)

var WebhookModule = fx.Options(
	fx.Provide(
		NewWebhookService,
		NewWebhookHandler,
	),
	fx.Invoke(
		RegisterWebhookRoutes,
	),
)
