package hello

import (
	"log"
	"sync"
	"time"
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
