package hello

import (
	"context"
	"log"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) HandleRoot(ctx context.Context, input *struct{}) (*HelloResponse, error) {
	log.Println("Handling root request")
	return &HelloResponse{Body: Model{
		Message: "Hello, World!",
	}}, nil
}
