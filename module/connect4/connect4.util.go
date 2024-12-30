package connect4

import (
	"connect4server/utils"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/gorilla/websocket"
)

const runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890123456789012345678901234567890"
const codeLength = 8

func generateGameCode() string {
	code := []byte{}

	for i := 0; i < codeLength; i++ {

		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		utils.Assert(err)

		code = append(code, runes[index.Int64()])
	}

	return string(code)
}

func (dto *MessageDto) Unmarshal(message []byte) {
	err := json.Unmarshal(message, dto)
	utils.Assert(err)
}

func (dto *MessageDto) Marshal() []byte {

	marshalledBytes, err := json.Marshal(dto)
	utils.Assert(err)

	return marshalledBytes
}

func GetConnectionId(conn *websocket.Conn) string {
	address := conn.RemoteAddr().String()
	parts := strings.Split(address, ":")
	return parts[0]
}
