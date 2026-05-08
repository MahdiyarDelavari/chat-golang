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

func (h *Hub) GetClients(userId int64) ([]*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	connections, ok := h.Clients[userId]
	if !ok || len(connections) == 0 {
		return nil, false
	}

	clients := make([]*Client, 0, len(connections))
	for client := range connections {
		clients = append(clients, client)
	}
	return clients, true
}

func (h *Hub) SendEventToUserIds(userIds []int64, eventType EventType, payload map[string]any) {
		for _, userId := range userIds {
			h.mu.RLock()
			connections, ok := h.Clients[userId]
			h.mu.RUnlock()
			if ok {
				for client := range connections {
					client.SendEvent(Event{
						EventType: eventType,
						Payload:   payload,
					})
				}
			}
		}
}