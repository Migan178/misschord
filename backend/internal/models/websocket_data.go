package models

import (
	"encoding/json"

	"github.com/Migan178/misschord-backend/internal/repository/ent"
)

type OPCode int

type WebSocketData struct {
	OP   *OPCode          `json:"op"`
	Data *json.RawMessage `json:"data,omitempty"`
	Type EventType        `json:"type,omitempty"`
}

const (
	OPCodeDispatch OPCode = iota
	OPCodeHeartBeat
	OPCodeHeartBeatACK
	OPCodeHello
	OPCodeIdentify
	OPCodeReady
	OPCodeError
)

type IdentifyData struct {
	Token string `json:"token"`
}

type HelloData struct {
	HeartbeatInterval int    `json:"heartbeatInterval"`
	Message           string `json:"message,omitempty"`
}

type ReadyData struct {
	User *ent.User `json:"user"`
}
