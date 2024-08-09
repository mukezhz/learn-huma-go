package domain

import (
    "github.com/mukezhz/learn-huma/domain/hello"
    "github.com/mukezhz/learn-huma/domain/middlewares"

    "go.uber.org/fx"
)

var Module = fx.Options(
	middlewares.Module,
	hello.Module,
)
