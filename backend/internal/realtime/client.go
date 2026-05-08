package realtime

import (
	"backend/internal/models"
	"sync"

	"github.com/coder/websocket"
)

type Client struct {
	User *models.User    `json:"user"`
	Conn *websocket.Conn `json:"-"`
	send chan Event      `json:"-"`
	once sync.Once           `json:"-"`
}

func NewClient(user *models.User, conn *websocket.Conn) *Client {
	return &Client{
		User: user,
		Conn: conn,
		send: make(chan Event, 512),
	}
}

func (c *Client) SendEvent(event Event) {
	select {
	case c.send <- event:
	default:
		// If the send channel is full, we can choose to drop the event or close the connection
		c.Close()
	}
}

func (c *Client) SendChannel() <-chan Event {
	return c.send
}

func (c *Client) Close() {
	c.once.Do(func() {
		if c.Conn != nil {
			c.Conn.Close(websocket.StatusNormalClosure, "Connection closed")
		}
		close(c.send)
	})
}
