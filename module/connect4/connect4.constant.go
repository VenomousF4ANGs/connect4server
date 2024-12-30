package connect4

const BOARD_COLUMNS = 7
const BOARD_ROWS = 6

type EnumGameState string

const (
	STATE_WAIT    EnumGameState = "WAIT"
	STATE_MATCHED EnumGameState = "MATCHED"
)

type EnumMessageType string

const (
	// exchange
	MESSAGE_PLAY EnumMessageType = "PLAY"
	MESSAGE_MOVE EnumMessageType = "MOVE"
	MESSAGE_TAP  EnumMessageType = "TAP"

	// common
	MESSAGE_GAME_END EnumMessageType = "GAME_END"
	MESSAGE_ERROR    EnumMessageType = "ERROR"

	// receive
	MESSAGE_MATCH_FIND   EnumMessageType = "MATCH_FIND"
	MESSAGE_MATCH_CANCEL EnumMessageType = "MATCH_CANCEL"

	// send ( no handlers )
	MESSAGE_DISCONNECTED EnumMessageType = "DISCONNECTED"
	MESSAGE_RECONNECTED  EnumMessageType = "RECONNECTED"
	MESSAGE_GAME_WAIT    EnumMessageType = "GAME_WAIT"
	MESSAGE_GAME_FOUND   EnumMessageType = "GAME_FOUND"
)

type EnumPlayerState string

const (
	PLAYER_WAITING      EnumPlayerState = "WAITING"
	PLAYER_PLAYING      EnumPlayerState = "PLAYING"
	PLAYER_DISCONNECTED EnumPlayerState = "DISCONNECTED"
)
