package realtime

import (
	"log"
	"sync"
)

type Hub struct {
	Clients map[int64]map[*Client]struct{}
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[int64]map[*Client]struct{}),
	}
}

func (h *Hub) broadcastToAll(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, connections := range h.Clients {
		for client := range connections {
			select {
			case client.send <- event:
			default:
				log.Printf("Client send channel full, closing connection for user %d", client.User.ID)
			}
		}
	}
}
