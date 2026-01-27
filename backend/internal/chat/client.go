package chat

import (
	"fmt"
	"time"

	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub           *Hub
	conn          *websocket.Conn
	user          *ent.User
	send          chan models.WebSocketData
	currentRoomID string
}

const (
	messageLimit = 8096
	clientWait   = 45 * time.Second
	writeWait    = 10 * time.Second
)

func NewClient(hub *Hub, user *ent.User, conn *websocket.Conn) *Client {
	client := &Client{hub, conn, user, make(chan models.WebSocketData, 512), ""}

	hub.register <- client

	return client
}

func (c *Client) Start() {
	go c.writePump()

	helloOP := models.OPCodeHello
	message := models.WebSocketData{OP: &helloOP}
	data := models.HelloData{
		HeartbeatInterval: int(clientWait) / 1000000,
	}

	if c.user == nil {
		data.Message = fmt.Sprintf("you need to authorization to send opcode %d", models.OPCodeIdentify)
	}

	c.parseDataAndSend(&message, data)

	if c.user != nil {
		readyOP := models.OPCodeReady
		c.parseDataAndSend(&models.WebSocketData{OP: &readyOP}, models.ReadyData{User: c.user})
	}

	c.readPump()
}
