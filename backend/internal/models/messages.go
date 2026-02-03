package models

import (
	"time"

	"github.com/Migan178/misschord-backend/internal/repository/ent"
)

type MessageCreateData struct {
	Message string       `json:"message" binding:"required"`
	Channel *ChannelData `json:"-"`
}

type MessageResponse struct {
	ID        int         `json:"id"`
	Author    *ent.User   `json:"author"`
	Message   string      `json:"message"`
	Channel   ChannelData `json:"channel"`
	CreatedAt time.Time   `json:"created_at"`
}
