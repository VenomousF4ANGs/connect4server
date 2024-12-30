package connect4

import (
	"connect4server/utils"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func StateService() *StateDto {
	return &StateDto{
		GameStore:       gameStore,
		MatchMaking:     matchMaking,
		ConnectionStore: connectionStore,
	}
}

func ProcessReconnection(conn *websocket.Conn) {
	gameStoreMutex.Lock()
	defer gameStoreMutex.Unlock()

	handleReConnection(conn)
}

func ProcessError(conn *websocket.Conn, err error) {
	if err != nil {
		gameStoreMutex.Lock()
		defer gameStoreMutex.Unlock()

		if _, ok := err.(*websocket.CloseError); ok {
			handleDisconnection(conn)
		} else {
			handleDisconnection(conn)
		}

		return
	}
}

func ProcessMessage(conn *websocket.Conn, messageType int, message []byte) {

	utils.ErrorIf(messageType != websocket.TextMessage)
	fmt.Println("Connect4 Server Message:", string(message))
	dto := MessageDto{}
	dto.Unmarshal(message)

	gameStoreMutex.Lock()
	defer gameStoreMutex.Unlock()

	switch dto.Type {
	case MESSAGE_PLAY:
		validatePlay(conn, &dto)
		updateState(&dto)
		updateSequence(conn, &dto)
		forwardToOtherPlayer(&dto)

	case MESSAGE_MOVE, MESSAGE_TAP:
		validatePlay(conn, &dto)
		updateSequence(conn, &dto)
		forwardToOtherPlayer(&dto)

	case MESSAGE_MATCH_FIND:
		validateMatchMaking(conn)
		startMatchMaking(conn, &dto)

	case MESSAGE_MATCH_CANCEL:
		validateMatchCancel(conn, &dto)
		cancelMatchMaking(conn, &dto)

	case MESSAGE_GAME_END, MESSAGE_ERROR:
		validatePlayer(conn, &dto)
		endGame(dto.GameCode)
	}
}

func updateSequence(conn *websocket.Conn, dto *MessageDto) {
	connectionId := GetConnectionId(conn)
	gameObject := getRunningGame(dto.GameCode)

	switch {
	case connectionId == gameObject.Player1.ConnectionId:
		utils.ErrorIf(dto.Sequence < gameObject.Player1.Sequence)
		gameObject.Player1.Sequence += 1

	case connectionId == gameObject.Player2.ConnectionId:
		utils.ErrorIf(dto.Sequence < gameObject.Player2.Sequence)
		gameObject.Player2.Sequence += 1

	default:
		utils.ErrorIf(true)
	}
}

func cancelMatchMaking(conn *websocket.Conn, messageDto *MessageDto) {
	// extra handling just in case
	gameObject := gameStore[messageDto.GameCode]
	if gameObject.State == STATE_MATCHED {
		notifyBothPlayersEnd(gameObject)
	}

	delete(gameStore, messageDto.GameCode)
	delete(matchMaking, messageDto.GameCode)
	delete(connectionStore, GetConnectionId(conn))
}

func validateMatchCancel(conn *websocket.Conn, messageDto *MessageDto) {
	connectionId := GetConnectionId(conn)

	gameCode, connectionExists := connectionStore[connectionId]
	utils.ErrorIf(!connectionExists)

	utils.ErrorIf(messageDto.GameCode != gameCode)
	gameObject, gameExists := gameStore[gameCode]
	utils.ErrorIf(!gameExists)

	utils.ErrorIf(gameObject.Player1.PlayerId != messageDto.PlayerId)
	utils.ErrorIf(gameObject.Player1.ConnectionId != connectionId)

}

func validateMatchMaking(conn *websocket.Conn) {
	connectionId := GetConnectionId(conn)

	gameCode, connectionExists := connectionStore[connectionId]
	if connectionExists {
		gameObject, gameExists := gameStore[gameCode]
		if gameExists && gameObject.State == STATE_MATCHED {

			notifyError(gameObject, "Duplicate Connection")
			notifyErrorAdhoc(conn, "Duplicate Connection")

			utils.ErrorIf(true)
		}
	}
}

func validatePlay(conn *websocket.Conn, dto *MessageDto) {
	gameObject := getRunningGame(dto.GameCode)

	utils.ErrorIf(gameObject == nil)
	utils.ErrorIf(gameObject.State != STATE_MATCHED)

	connectionId := GetConnectionId(conn)

	switch {
	case connectionId == gameObject.Player1.ConnectionId:
		utils.ErrorIf(gameObject.Player1.PlayerId != dto.PlayerId)
		utils.ErrorIf(gameObject.Turn != 1)

	case connectionId == gameObject.Player2.ConnectionId:
		utils.ErrorIf(gameObject.Player2.PlayerId != dto.PlayerId)
		utils.ErrorIf(gameObject.Turn != 2)
	default:
		utils.ErrorIf(true)
	}

	utils.ErrorIf(dto.Column < 0 || dto.Column >= BOARD_COLUMNS)
	utils.ErrorIf(gameObject.BoardGrounds[dto.Column] >= (BOARD_ROWS - 1))

}

func validatePlayer(conn *websocket.Conn, dto *MessageDto) {
	gameObject := getRunningGame(dto.GameCode)

	utils.ErrorIf(gameObject == nil)
	utils.ErrorIf(gameObject.State != STATE_MATCHED)

	connectionId := GetConnectionId(conn)

	switch {
	case connectionId == gameObject.Player1.ConnectionId:
		utils.ErrorIf(gameObject.Player1.PlayerId != dto.PlayerId)

	case connectionId == gameObject.Player2.ConnectionId:
		utils.ErrorIf(gameObject.Player2.PlayerId != dto.PlayerId)

	default:
		utils.ErrorIf(true)
	}
}

func endGame(gameCode string) {
	gameObject, gameExists := gameStore[gameCode]
	if !gameExists {
		return
	}

	notifyBothPlayersEnd(gameObject)

	delete(gameStore, gameCode)
	delete(matchMaking, gameCode)

	if gameObject.Player1 != nil {
		delete(connectionStore, gameObject.Player1.ConnectionId)
	}

	if gameObject.Player2 != nil {
		delete(connectionStore, gameObject.Player2.ConnectionId)
	}
}

func forwardToOtherPlayer(dto *MessageDto) {
	gameObject := getRunningGame(dto.GameCode)
	utils.ErrorIf(gameObject == nil)

	switch dto.PlayerId {

	case gameObject.Player1.PlayerId:
		dto.PlayerId = gameObject.Player2.PlayerId
		gameObject.Player2.Connection.WriteMessage(websocket.TextMessage, dto.Marshal())

	case gameObject.Player2.PlayerId:
		dto.PlayerId = gameObject.Player1.PlayerId
		gameObject.Player1.Connection.WriteMessage(websocket.TextMessage, dto.Marshal())

	}
}

func handleReConnection(conn *websocket.Conn) bool {
	connectionId := GetConnectionId(conn)
	gameCode, connectionExists := connectionStore[connectionId]

	if !connectionExists {
		// notifyErrorAdhoc(conn, "Test Error")
		return false
	}

	gameObject, gameExists := gameStore[gameCode]

	if !gameExists {
		delete(connectionStore, connectionId)
		notifyErrorAdhoc(conn, "Game Deleted")
		return false
	}

	switch {
	case gameObject.Player1 != nil && connectionId == gameObject.Player1.ConnectionId:
		if gameObject.State == STATE_MATCHED {
			gameObject.Player1.State = PLAYER_PLAYING
			gameObject.Player2.sendReconnection("Opponent reconnected")
		} else {
			gameObject.Player1.State = PLAYER_WAITING
		}

		return true

	case gameObject.Player2 != nil && connectionId == gameObject.Player2.ConnectionId:
		if gameObject.State == STATE_MATCHED {
			gameObject.Player2.State = PLAYER_PLAYING
			gameObject.Player1.sendReconnection("Opponent reconnected")
		} else {
			gameObject.Player2.State = PLAYER_WAITING
		}

		return true
	}

	return false

}

func handleDisconnection(conn *websocket.Conn) {
	connectionId := GetConnectionId(conn)
	gameCode, connectionExists := connectionStore[connectionId]

	if !connectionExists {
		return
	}

	gameObject, gameExists := gameStore[gameCode]

	if !gameExists {
		delete(connectionStore, connectionId)
		return
	}

	switch {
	case gameObject.Player1 != nil && connectionId == gameObject.Player1.ConnectionId:

		if gameObject.Player1.State != PLAYER_DISCONNECTED {
			gameObject.Player1.State = PLAYER_DISCONNECTED
			time.AfterFunc(10*time.Second, func() {
				if gameObject.Player1.State == PLAYER_DISCONNECTED {
					endGame(gameCode)
				}
			})
			if gameObject.State == STATE_MATCHED {
				gameObject.Player2.sendDisconnection("Opponent disconnected")
			}
		}

	case gameObject.Player2 != nil && connectionId == gameObject.Player2.ConnectionId:

		if gameObject.Player1.State != PLAYER_DISCONNECTED {
			gameObject.Player2.State = PLAYER_DISCONNECTED
			time.AfterFunc(10*time.Second, func() {
				if gameObject.Player2.State == PLAYER_DISCONNECTED {
					endGame(gameCode)
				}
			})
			if gameObject.State == STATE_MATCHED {
				gameObject.Player1.sendDisconnection("Opponent disconnected")
			}
		}

	}
}

func startMatchMaking(conn *websocket.Conn, dto *MessageDto) {

	gameObject := getFreeMatch(conn, dto)

	switch {

	case gameObject == nil:
		gameObject = addToMatchMaking(conn, dto)
		askPlayer1ToWait(gameObject)

	case gameObject.Player1.ConnectionId == GetConnectionId(conn):
		askPlayer1ToWait(gameObject)

	default:
		gameObject.Turn = 1
		notifyBothPlayersFound(gameObject)

	}
}

func updateState(dto *MessageDto) {
	gameObject := getRunningGame(dto.GameCode)

	gameObject.BoardGrounds[dto.Column] += 1
	gameObject.Turn = (gameObject.Turn % 2) + 1 // toggle between 1 <-> 2
}
