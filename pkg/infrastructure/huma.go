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
	api.OpenAPI().Info.Contact = &huma.Contact{
		Name:  "Mukesh Kumar Chaudhary",
		Email: "mukezhz@duck.com",
		URL:   "https://mukesh.name.np",
	}
	api.OpenAPI().Info.License = &huma.License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}
	api.OpenAPI().Info.Description = "This is a demo API for learning Huma."
	return &HumaRouter{
		api,
	}
}
