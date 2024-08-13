package hello

import (
	"context"
	"log"

	"github.com/mukezhz/learn-huma/pkg/infrastructure"
)

type Controller struct {
	service *Service
	hello   *infrastructure.Hello
}

func NewController(
	service *Service,
	hello *infrastructure.Hello,
) *Controller {
	return &Controller{service, hello}
}

func (ctrl *Controller) HandleRoot(ctx context.Context, input *struct{}) (*HelloResponse, error) {
	log.Println("Handling root request")

	return &HelloResponse{Body: Model{
		Message: "Hello, World!",
	}}, nil
}
