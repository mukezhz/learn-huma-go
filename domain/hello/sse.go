package hello

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type Message struct {
	UserID       string `json:"user_id"`
	Message      string `json:"message"`
	BroadCastAll bool   `json:"all"`
}

type ClientChan chan string

type SSEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
	ID    string `json:"id"`
	Retry int    `json:"retry"`
}

type SSEInput struct {
	Query         string `query:"q"`
	Authorization string `header:"Authorization"`
}

type EventStream struct {
	Message       chan Message
	NewClients    chan ClientChan
	ClosedClients chan ClientChan
	TotalClients  map[ClientChan]string
	UserIDChanMap map[string]ClientChan
	Ticker        *time.Ticker
	mu            sync.RWMutex
}

func NewEventStream() *EventStream {
	event := &EventStream{
		Message:       make(chan Message, 10),
		NewClients:    make(chan ClientChan),
		ClosedClients: make(chan ClientChan),
		TotalClients:  make(map[ClientChan]string),
		UserIDChanMap: make(map[string]ClientChan),
		Ticker:        time.NewTicker(time.Second * 10),
	}

	go event.listen()
	return event
}

// Listen to the ticker and send heartbeat to the client
func (e *EventStream) listen() {
	for {
		select {
		case newClient := <-e.NewClients:
			e.mu.Lock()
			e.TotalClients[newClient] = <-newClient
			e.mu.Unlock()
			log.Printf("Client added. %d registered clients", len(e.TotalClients))

		case closeClient := <-e.ClosedClients:
			e.mu.Lock()
			log.Printf("Removing client!!!")
			delete(e.TotalClients, closeClient)
			close(closeClient)
			log.Printf("Removed client. %d registered clients", len(e.TotalClients))
			e.mu.Unlock()

		case eventMsg := <-e.Message:
			log.Println("Received message!!!", eventMsg)
			e.mu.RLock()
			for clientMessageChan, userID := range e.TotalClients {
				log.Println("Sending message!!!", userID, userID == eventMsg.UserID, eventMsg.BroadCastAll)
				if eventMsg.BroadCastAll || userID == eventMsg.UserID {
					clientMessageChan <- eventMsg.Message
				}
			}
			e.mu.RUnlock()

		case t := <-e.Ticker.C:
			e.Message <- Message{
				Message:      t.Format(time.RFC3339),
				BroadCastAll: true,
			}
		}
	}
}

// SIMPLE SSE
type Client struct {
	ID      string
	Channel chan string
}

type SSE struct {
	clients     map[string]*Client
	clientMutex sync.RWMutex
}

func NewSSE() *SSE {
	return &SSE{
		clients: make(map[string]*Client),
	}
}

func (s *SSE) AddClient(id string, channel chan string) chan string {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	client := Client{
		ID:      id,
		Channel: make(chan string),
	}
	s.clients[id] = &client
	return client.Channel
}

func (s *SSE) RemoveClient(id string) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()
	if client, ok := s.clients[id]; ok {
		close(client.Channel) // Close the channel when the client disconnects
		delete(s.clients, id)
	}
}

func (s *SSE) SendMessageToUser(userIDs []string, message string) {
	s.clientMutex.RLock()
	defer s.clientMutex.RUnlock()
	for _, id := range userIDs {
		if client, ok := s.clients[id]; ok {
			client.Channel <- message
		}
	}
}

func (s *SSE) BroadcastMessage(message string) {
	s.clientMutex.RLock()
	defer s.clientMutex.RUnlock()
	for _, client := range s.clients {
		client.Channel <- message
	}
}

func (s *SSE) SendMessage(ctx huma.Context, message string) {
	_, err := fmt.Fprintf(ctx.BodyWriter(), "data: %s\n\n", message)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
	if flusher, ok := ctx.BodyWriter().(http.Flusher); ok {
		flusher.Flush()
	} else {
		huma.Error503ServiceUnavailable("Streaming unsupported!")
	}
}
