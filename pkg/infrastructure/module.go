package infrastructure

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Module(
	"infrastructure",
	fx.Options(
		fx.Provide(NewRouter),
		fx.Provide(NewHumaRouter),
		fx.Provide(NewHello),
	),
)
