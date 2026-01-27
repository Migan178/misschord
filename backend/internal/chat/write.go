package chat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Migan178/misschord-backend/internal/models"
)

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		c.conn.SetWriteDeadline(time.Now().Add(writeWait))

		if err := c.conn.WriteJSON(message); err != nil {
			fmt.Printf("write to client err: %v\n", err)
			return
		}
	}
}

func (c *Client) parseDataAndSend(message *models.WebSocketData, data any) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return
	}

	message.Data = (*json.RawMessage)(&dataByte)
	c.safeSend(message)
}

func (c *Client) safeSend(message *models.WebSocketData) {
	select {
	case c.send <- *message:
		return
	case <-time.After(writeWait):
		c.hub.unregister <- c
	}
}
