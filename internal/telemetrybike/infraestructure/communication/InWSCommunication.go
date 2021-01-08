package communication

import (
	"app/internal/telemetrybike/domain"

	"github.com/gorilla/websocket"
)

type communication struct {
	conn *websocket.Conn
}

// NewWSCommunication return a communication instance to send data over gorilla websocket package
func NewWSCommunication(c *websocket.Conn) domain.TelemetryBikeCommunication {
	return &communication{
		conn: c,
	}
}

func (c *communication) Send(content interface{}) error {
	err := c.conn.WriteJSON(content)
	if err != nil {
		return err
	}
	return nil
}
