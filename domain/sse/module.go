package sse

import (
	"go.uber.org/fx"
)

var SseModule = fx.Options(
	fx.Provide(
		NewBroadcaster,
	),
	fx.Invoke(
		RegisterSSERoutes,
	),
)
