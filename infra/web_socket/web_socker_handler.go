package web_socket

import (
	"connect4server/module/connect4"
	"connect4server/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsServiceFunction func(*websocket.Conn, int, []byte)
type WsErrorFunction func(*websocket.Conn, error)
type WsReconnectionFunction func(*websocket.Conn)

type WsHandler struct {
	ServiceFunction      WsServiceFunction
	ReconnectionFunction WsReconnectionFunction
	ErrorFunction        WsErrorFunction
}

func (handler *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	utils.Assert(err)
	defer c.Close()

	fmt.Println("Conection Upgraded Id: ", connect4.GetConnectionId(c))
	handler.ReconnectionFunction(c)

	for {
		messageType, message, err := c.ReadMessage()

		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				handler.ErrorFunction(c, err)
			} else {
				handler.ErrorFunction(c, err)
			}
		} else {
			handler.ServiceFunction(c, messageType, message)
		}

	}
}

func httpGetHandler[DTO any](target func() *DTO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := target()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(dto)
		utils.Assert(err)
	}
}
