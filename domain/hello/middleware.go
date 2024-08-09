package hello

import (
	"github.com/danielgtaylor/huma/v2"
)

func SSEHeaderMiddleware(ctx huma.Context, next func(huma.Context)) {
	ctx.SetHeader("Content-Type", "text/event-stream")
	ctx.SetHeader("Cache-Control", "no-cache")
	ctx.SetHeader("Connection", "keep-alive")
	ctx.SetHeader("Transfer-Encoding", "chunked")

	next(ctx)
}

// Middleware for managing SSE connection lifecycle
func (e *EventStream) serveSSE(ctx huma.Context, next func(huma.Context)) {
	userID := ctx.Query("q")
	clientChan := make(ClientChan, 10)
	clientChan <- userID
	e.UserIDChanMap[userID] = clientChan
	e.NewClients <- clientChan
	defer func() {
		e.ClosedClients <- clientChan
	}()

	ctx = huma.WithValue(ctx, "clientChan", clientChan)

	next(ctx)
}
