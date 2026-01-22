package instagram

import (
	"go.uber.org/fx"
)

var InstagramModule = fx.Options(
	fx.Provide(
		NewInstagramHandler,
	),
	fx.Invoke(
		RegisterInstagramRoutes,
	),
)
