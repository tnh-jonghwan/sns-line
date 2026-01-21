package eventHub

import (
	"go.uber.org/fx"
)

var EventHubModule = fx.Options(
	fx.Provide(
		NewEventHub,
	),
	fx.Invoke(
		RegistereventHubRoutes,
	),
)
