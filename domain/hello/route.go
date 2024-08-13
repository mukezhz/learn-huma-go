package hello

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mukezhz/learn-huma/pkg/infrastructure"
)

type Route struct {
	controller *Controller
	api        *infrastructure.HumaRouter
	router     *infrastructure.Router
	stream     *EventStream
}

func NewRoute(
	controller *Controller,
	router *infrastructure.Router,
	api *infrastructure.HumaRouter,
	stream *EventStream,
) *Route {
	return &Route{
		controller: controller,
		api: &infrastructure.HumaRouter{
			API: api,
		},
		router: router,
		stream: stream,
	}
}

func RegisterRoute(r *Route) {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			time.Sleep(time.Second * 1)
			log.Println("Sending Current Time")
			now := time.Now().Format("2006-01-02 15:04:05")
			currentTime := fmt.Sprintf("The Current Time Is %v", now)
			r.stream.Message <- Message{
				Message:      currentTime,
				UserID:       "123",
				BroadCastAll: false,
			}
		}
	}()

	r.stream.Ticker = ticker

	huma.Register(r.api, huma.Operation{
		OperationID: "hello",
		Method:      http.MethodGet,
		Path:        "/hello",
		Summary:     "Hello, world!",
	}, r.controller.HandleRoot,
	)

	huma.Get(r.api, "/hello/demo", r.controller.HandleRoot)
	op := getSSEOperation(r.api.API, "/sse-gin", map[string]any{
		"message":      DefaultMessage{},
		"userCreate":   UserCreatedEvent{},
		"mailReceived": MailReceivedEvent{},
	})
	// op.Middlewares = huma.Middlewares{SSEHeaderMiddleware, r.stream.serveSSE}
	// huma.Register(r.api, op, func(ctx context.Context, input *SSEInput) (*huma.StreamResponse, error) {
	// 	return &huma.StreamResponse{
	// 		Body: func(ctx huma.Context) {
	// 			clientChan, ok := ctx.Context().Value("clientChan").(ClientChan)
	// 			if !ok {
	// 				log.Println("clientChan not found")
	// 				return
	// 			}

	// 			// Use errgroup to manage the goroutine
	// 			g, gctx := errgroup.WithContext(ctx.Context())

	// 			g.Go(func() error {
	// 				for {
	// 					select {
	// 					case <-gctx.Done():
	// 						return gctx.Err() // Return context error
	// 					case msg := <-clientChan:
	// 						send(ctx, SSEvent{
	// 							Event: "message",
	// 							Data:  msg,
	// 						})
	// 					}
	// 				}
	// 			})

	// 			// Wait for the errgroup to finish
	// 			if err := g.Wait(); err != nil {
	// 				log.Println("Error occurred: ", err)
	// 			}
	// 		},
	// 	}, nil
	// })

	op = getSSEOperation(r.api.API, "/sse", map[string]any{
		"message":      DefaultMessage{},
		"userCreate":   UserCreatedEvent{},
		"mailReceived": MailReceivedEvent{},
	})
	op.Middlewares = huma.Middlewares{SSEHeaderMiddleware, r.stream.serveSSE}
	huma.Register(r.api, op, func(ctx context.Context, input *SSEInput) (*huma.StreamResponse, error) {
		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				clientID := fmt.Sprintf("%d", time.Now().UnixNano())
				messageChannel := make(chan string)

				sse := NewSSE()
				sse.AddClient(clientID, messageChannel)
				defer sse.RemoveClient(clientID)

				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()

				go func() {
					for msg := range ticker.C {
						sse.SendMessage(ctx, "The Current Time Is "+msg.Format("2006-01-02 15:04:05"))
					}
				}()

				go func() {
					for msg := range messageChannel {
						sse.SendMessage(ctx, msg)
					}
				}()

				<-ctx.Context().Done()
				fmt.Println("Client disconnected")
			},
		}, nil
	})

}

var WriteTimeout = 5 * time.Second

type DefaultMessage struct {
	Message string `json:"message"`
}

type UserCreatedEvent struct {
	UserID string `json:"userId"`
}

type MailReceivedEvent struct {
	UserID string `json:"userId"`
}
