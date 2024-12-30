package connect4

import (
	"connect4server/utils"
	"fmt"
	"slices"

	"github.com/gorilla/websocket"
)

func getFreeMatch(connection *websocket.Conn, dto *MessageDto) *GameObject {

	for gameCode := range matchMaking {

		gameObject, gameExists := gameStore[gameCode]
		utils.ErrorIf(!gameExists)
		utils.ErrorIf(gameObject.State != STATE_WAIT)

		switch {

		case gameObject.Player1.ConnectionId == GetConnectionId(connection):
			gameObject.Player1 = newPlayerObject(connection, dto.PlayerId, gameCode)
			gameObject.Player1.Sequence = dto.Sequence

			return gameObject

		case gameObject.Player1.State == PLAYER_DISCONNECTED:
			continue

		default:
			gameObject.Player2 = newPlayerObject(connection, dto.PlayerId, gameCode)
			gameObject.Player2.Sequence = dto.Sequence

			gameObject.State = STATE_MATCHED
			gameObject.Player2.State = PLAYER_PLAYING
			gameObject.Player1.State = PLAYER_PLAYING

			delete(matchMaking, gameCode)
			return gameObject
		}
	}

	return nil
}

func addToMatchMaking(player1 *websocket.Conn, dto *MessageDto) *GameObject {
	gameObject := newGameObject(player1, dto)

	gameStore[gameObject.GameCode] = gameObject
	matchMaking[gameObject.GameCode] = true

	return gameObject
}

func getRunningGame(code string) *GameObject {
	gameObject, exists := gameStore[code]
	if exists && gameObject.State == STATE_MATCHED {
		return gameObject
	}

	return nil
}

func newGameObject(conn *websocket.Conn, dto *MessageDto) *GameObject {
	newGameCode := generateGameCode()
	gameObject := GameObject{
		GameCode:     newGameCode,
		State:        STATE_WAIT,
		BoardGrounds: slices.Repeat([]int{0}, BOARD_COLUMNS),
	}

	gameObject.Player1 = newPlayerObject(conn, dto.PlayerId, newGameCode)
	gameObject.Player1.Sequence = dto.Sequence

	connectionStore[gameObject.Player1.ConnectionId] = gameObject.GameCode

	return &gameObject
}

func newPlayerObject(conn *websocket.Conn, id string, gameCode string) *PlayerObject {
	playerObject := PlayerObject{
		PlayerId:     id,
		ConnectionId: GetConnectionId(conn),
		Connection:   conn,
		Sequence:     1,
		State:        PLAYER_WAITING,
	}

	connectionStore[playerObject.ConnectionId] = gameCode

	return &playerObject
}

func (player *PlayerObject) sendData(dto *MessageDto, updateSeq bool) {

	if player.State != PLAYER_DISCONNECTED {
		if updateSeq {
			player.Sequence += 1
		}

		dto.Sequence = player.Sequence
		textData := dto.Marshal()
		fmt.Println("Send Data:", string(textData))
		player.Connection.WriteMessage(websocket.TextMessage, textData)
	}
}

func notifyErrorAdhoc(conn *websocket.Conn, message string) {
	dto := MessageDto{
		Type:    MESSAGE_ERROR,
		Message: message,
	}

	textData := dto.Marshal()
	fmt.Println("Send Data:", string(textData))
	conn.WriteMessage(websocket.TextMessage, textData)
}

func (player *PlayerObject) sendDisconnection(message string) {
	dto := MessageDto{
		Type:     MESSAGE_DISCONNECTED,
		Message:  message,
		Sequence: player.Sequence,
	}

	if player.State != PLAYER_DISCONNECTED {
		textData := dto.Marshal()
		fmt.Println("Send Data:", string(textData))
		player.Connection.WriteMessage(websocket.TextMessage, textData)
	}
}

func (player *PlayerObject) sendReconnection(message string) {
	dto := MessageDto{
		Type:     MESSAGE_RECONNECTED,
		Message:  message,
		Sequence: player.Sequence,
	}

	if player.State != PLAYER_DISCONNECTED {
		textData := dto.Marshal()
		fmt.Println("Send Data:", string(textData))
		player.Connection.WriteMessage(websocket.TextMessage, textData)
	}
}
