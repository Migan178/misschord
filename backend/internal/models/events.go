package models

import (
	"time"

	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/room"
)

type EventType string

const (
	EventTypeMessageCreate EventType = "MESSAGE_CREATE"
	EventTypeChannelJoin   EventType = "CHANNEL_JOIN"
	EventTypeChannelLeave  EventType = "CHANNEL_LEAVE"
)

type MessageCreateEvent struct {
	ID        *int        `json:"id"`
	Author    *ent.User   `json:"author"`
	Message   string      `json:"message"`
	Channel   ChannelData `json:"channel"`
	CreatedAt *time.Time  `json:"createdAt"`
}

func (m *MessageCreateEvent) GetInternalRoomID() int {
	return m.Channel.ID
}

type ChannelData struct {
	ID       int           `json:"id"`
	RoomType room.RoomType `json:"type"`
}
