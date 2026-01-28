package chat

import (
	"encoding/json"

	"github.com/Migan178/misschord-backend/internal/models"
	jwt "github.com/appleboy/gin-jwt/v3"
)

type joinOrLeaveData struct {
	client *Client
	roomID string
}

type Hub struct {
	clients        map[*Client]bool
	authMiddleware *jwt.GinJWTMiddleware
	rooms          map[string]map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	join           chan joinOrLeaveData
	leave          chan joinOrLeaveData
	broadcast      chan broadcastData
}

type broadcastData struct {
	message *models.WebSocketData
	roomID  string
}

type dataToBroadcast interface {
	GetInternalRoomID() string
}

func NewHub(authMiddleware *jwt.GinJWTMiddleware) *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		authMiddleware: authMiddleware,
		rooms:          make(map[string]map[*Client]bool),
		register:       make(chan *Client, 256),
		unregister:     make(chan *Client, 256),
		join:           make(chan joinOrLeaveData, 256),
		leave:          make(chan joinOrLeaveData, 256),
		broadcast:      make(chan broadcastData, 512),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.disconnect()
			}
		case data := <-h.broadcast:
			switch data.message.Type {
			case models.EventTypeMessageCreate:
				if clients, ok := h.rooms[data.roomID]; ok {
					for client := range clients {
						client.safeSend(data.message)
					}
				}
			}
		case data := <-h.join:
			if room, ok := h.rooms[data.roomID]; ok {
				if _, ok := room[data.client]; ok {
					continue
				}

				room[data.client] = true
				continue
			}

			h.rooms[data.roomID] = make(map[*Client]bool)
			h.rooms[data.roomID][data.client] = true
		case data := <-h.leave:
			if room, ok := h.rooms[data.roomID]; ok {
				delete(room, data.client)
				data.client.currentRoomID = ""

				if len(room) == 0 {
					delete(h.rooms, data.roomID)
				}
			}
		}
	}
}

func (h *Hub) parseDataAndBroadcast(message *models.WebSocketData, data dataToBroadcast) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return
	}

	message.Data = (*json.RawMessage)(&dataByte)

	h.broadcast <- broadcastData{message, data.GetInternalRoomID()}
}
