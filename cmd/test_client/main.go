package main

import (
	"connect4server/module/connect4"
	"connect4server/utils"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:1234/connect4", nil)
	utils.Assert(err)
	defer c.Close()

	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	// read and display message
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}

			fmt.Printf("Received Response: %s\n", message)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	dto := connect4.MessageDto{
		Type:     connect4.MESSAGE_GAME_FOUND,
		GameCode: "abcd",
		PlayerId: "",
		Column:   0,
		Message:  "abcd",
	}
	err = c.WriteMessage(websocket.TextMessage, dto.Marshal())
	if err != nil {
		fmt.Println("write:", err)
		return
	}

}
