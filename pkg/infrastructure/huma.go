package infrastructure

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
)

type HumaRouter struct {
	huma.API
}

func NewHumaRouter(router *Router) *HumaRouter {
	api := humagin.New(router.Engine, huma.DefaultConfig("My API", "v0.0.1"))
	return &HumaRouter{
		api,
	}
}
