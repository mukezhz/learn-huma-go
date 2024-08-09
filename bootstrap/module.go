package bootstrap

import (
	"github.com/mukezhz/learn-huma/domain"
	"github.com/mukezhz/learn-huma/migrations"
	"github.com/mukezhz/learn-huma/pkg"
	"github.com/mukezhz/learn-huma/seeds"

	"go.uber.org/fx"
)

var CommonModules = fx.Module("common",
	fx.Options(
		pkg.Module,
		seeds.Module,
		migrations.Module,
		domain.Module,
	),
)
