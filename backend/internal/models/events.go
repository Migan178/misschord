package models

import (
	"time"

	"github.com/Migan178/misschord-backend/internal/repository/ent"
)

type EventType string

const (
	EventTypeMessageCreate EventType = "MESSAGE_CREATE"
	EventTypeChannelJoin   EventType = "CHANNEL_JOIN"
	EventTypeChannelLeave  EventType = "CHANNEL_LEAVE"
)

type ChannelType string

const (
	ChannelTypeDM ChannelType = "DM"
)

type MessageCreateEvent struct {
	ID        *int        `json:"id"`
	Author    *ent.User   `json:"author"`
	Message   string      `json:"message"`
	Channel   ChannelData `json:"channel"`
	CreatedAt *time.Time  `json:"createdAt"`
}

func (m *MessageCreateEvent) GetInternalRoomID() string {
	return m.Channel.InternalID
}

type ChannelData struct {
	ID          string      `json:"id"`
	ChannelType ChannelType `json:"type"`
	InternalID  string      `json:"-"`
}
