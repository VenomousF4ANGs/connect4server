package connect4

import "sync"

var gameStoreMutex sync.Mutex

var (
	gameStore       = map[string]*GameObject{}
	matchMaking     = map[string]bool{}
	connectionStore = map[string]string{}
)

type StateDto struct {
	GameStore       map[string]*GameObject
	MatchMaking     map[string]bool
	ConnectionStore map[string]string
}
