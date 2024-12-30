package connect4

import (
	"github.com/gorilla/websocket"
)

type GameObject struct {
	GameCode     string
	State        EnumGameState
	Turn         int
	BoardGrounds []int

	Player1 *PlayerObject
	Player2 *PlayerObject
}

type PlayerObject struct {
	PlayerId     string
	ConnectionId string
	Connection   *websocket.Conn
	State        EnumPlayerState
	Sequence     uint64
}

type MessageDto struct {
	Type     EnumMessageType `json:"type"`
	GameCode string          `json:"gameCode"`
	PlayerId string          `json:"playerId"`
	Message  string          `json:"message"`
	Column   int             `json:"column"`
	Sequence uint64          `json:"sequence"`
}
