package hello

import (
	"go.uber.org/fx"
)

var Module = fx.Module("hello",
	fx.Options(
		fx.Provide(
			NewService,
			NewController,
			NewRepository,
			NewRoute,
			NewEventStream,
		),
		fx.Invoke(RegisterRoute),
	),
)
