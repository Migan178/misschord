package models

import (
	"github.com/Migan178/misschord-backend/internal/repository/ent/room"
)

type EventType string

const (
	EventTypeMessageCreate EventType = "MESSAGE_CREATE"
	EventTypeChannelJoin   EventType = "CHANNEL_JOIN"
	EventTypeChannelLeave  EventType = "CHANNEL_LEAVE"
)

func (m *MessageResponse) GetInternalRoomID() int {
	return m.Channel.ID
}

type ChannelData struct {
	ID       int           `json:"id"`
	RoomType room.RoomType `json:"type"`
}
