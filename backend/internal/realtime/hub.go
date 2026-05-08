package realtime

import (
	"backend/internal/models"
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

func (h *Hub) SendEventToUserIds(userIds []int64, sendId int64, eventType EventType, payload map[string]any) {
	for _, id := range userIds {
		h.mu.RLock()
		conns, ok := h.Clients[id]
		h.mu.RUnlock()

		if !ok {
			continue
		}

		for c := range conns {
			c.SendEvent(Event{
				EventType: eventType,
				Payload:   payload,
			})
		}
	}
}

func (h *Hub) RegisterClientConnection(client *Client) {
	h.mu.Lock()
	connections, ok := h.Clients[client.User.ID]
	if !ok {
		connections = make(map[*Client]struct{})
		h.Clients[client.User.ID] = connections
	}
	connections[client] = struct{}{}
	firstConnection := len(connections) == 1
	h.mu.Unlock()

	if firstConnection {
		h.broadcastToAll(Event{
			EventType: EventUserOnline,
			Payload: map[string]any{
				"user_id": client.User.ID,
				"name":    client.User.Name,
				"email":   client.User.Email,
			},
		})

		go func() {
			privates, err := models.GetPrivatesForUser(client.User.ID)
			if err != nil {
				log.Printf("Error fetching privates for user %d: %v", client.User.ID, err)
				return
			}
			for _, private := range privates {
					msgs, err := models.GetUndeliveredMessagesByPrivateID(private.ID)
					if err != nil {
						log.Printf("Error fetching undelivered messages for private %d: %v", private.ID, err)
						continue
					}
					for _, msg := range msgs {
						if msg.FromID == client.User.ID {
							continue
						}
						h.SendEventToUserIds([]int64{msg.FromID}, client.User.ID,EventDelivered , map[string]any{
							"message_id": msg.ID,
							"to_id": client.User.ID,
						})
					}

			}
		}()
	}
}