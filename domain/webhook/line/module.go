package line

import (
	"go.uber.org/fx"
)

var LineWebhookModule = fx.Options(
	fx.Provide(
		NewWebhookService,
		NewWebhookHandler,
	),
	fx.Invoke(
		RegisterWebhookRoutes,
	),
)
