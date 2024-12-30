package connect4

func notifyBothPlayersFound(gameObject *GameObject) {
	dto := MessageDto{
		GameCode: gameObject.GameCode,
		Type:     MESSAGE_GAME_FOUND,
		Message:  "Game found on Server",
	}

	dto.PlayerId = gameObject.Player1.PlayerId
	dto.Column = 1
	gameObject.Player1.sendData(&dto, true)

	dto.PlayerId = gameObject.Player2.PlayerId
	dto.Column = 2
	gameObject.Player2.sendData(&dto, true)
}

func notifyBothPlayersEnd(gameObject *GameObject) {
	dto := MessageDto{
		GameCode: gameObject.GameCode,
		Type:     MESSAGE_GAME_END,
		Message:  "Game ended",
	}

	if gameObject.Player1 != nil {
		dto.PlayerId = gameObject.Player1.PlayerId
		gameObject.Player1.sendData(&dto, false)
	}

	if gameObject.Player2 != nil {
		dto.PlayerId = gameObject.Player2.PlayerId
		gameObject.Player2.sendData(&dto, false)
	}
}

func askPlayer1ToWait(gameObject *GameObject) {
	dto := MessageDto{
		GameCode: gameObject.GameCode,
		Type:     MESSAGE_GAME_WAIT,
		Message:  "Wait for Match",
	}

	dto.PlayerId = gameObject.Player1.PlayerId
	gameObject.Player1.sendData(&dto, true)
}

func notifyError(gameObject *GameObject, message string) {
	dto := MessageDto{
		GameCode: gameObject.GameCode,
		Type:     MESSAGE_ERROR,
		Message:  message,
	}

	if gameObject.Player1 != nil {
		dto.PlayerId = gameObject.Player1.PlayerId
		gameObject.Player1.sendData(&dto, false)
	}

	if gameObject.Player2 != nil {
		dto.PlayerId = gameObject.Player2.PlayerId
		gameObject.Player2.sendData(&dto, false)
	}
}
