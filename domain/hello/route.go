package hello

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/mukezhz/learn-huma/pkg/infrastructure"
)

type Route struct {
	controller *Controller
	api        *infrastructure.HumaRouter
}

func NewRoute(
	controller *Controller,
	router *infrastructure.Router,
	api *infrastructure.HumaRouter,
) *Route {
	return &Route{
		controller: controller,
		api: &infrastructure.HumaRouter{
			API: api,
		},
	}
}

func RegisterRoute(r *Route) {
	huma.Get(r.api, "/hello/demo", r.controller.HandleRoot)
}
