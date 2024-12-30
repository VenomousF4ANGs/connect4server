package echo

import (
	"connect4server/module/connect4"
	"fmt"

	"github.com/gorilla/websocket"
)

func EchoService(conn *websocket.Conn, mt int, message []byte) {

	conn.WriteMessage(mt, message)
}

func EchoError(conn *websocket.Conn, err error) {
	if err != nil {
		if _, ok := err.(*websocket.CloseError); ok {
			// fmt.Println("Echo Server Close Error:")
		} else {
			fmt.Println("Echo Server Unknown Error:", err)
		}

		return
	}
}

func EchoReconnection(conn *websocket.Conn) {
	fmt.Println("Echo Server Reconnect")
	conn.WriteMessage(websocket.TextMessage, []byte(connect4.GetConnectionId(conn)))
}
