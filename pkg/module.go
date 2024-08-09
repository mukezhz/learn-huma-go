package pkg

import (
	"github.com/mukezhz/learn-huma/pkg/framework"
	"github.com/mukezhz/learn-huma/pkg/infrastructure"
	"github.com/mukezhz/learn-huma/pkg/middlewares"
	"github.com/mukezhz/learn-huma/pkg/services"

	"go.uber.org/fx"
)

var Module = fx.Module("pkg",
	framework.Module,
	services.Module,
	middlewares.Module,
	infrastructure.Module,
)
