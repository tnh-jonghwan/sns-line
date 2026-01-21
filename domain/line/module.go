package line

import (
	"go.uber.org/fx"
)

var LineModule = fx.Options(
	fx.Provide(
		NewLineClient,
	),
)
