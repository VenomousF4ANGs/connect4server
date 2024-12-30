package web_socket

import (
	"connect4server/module/connect4"
	"connect4server/module/echo"
	"connect4server/utils"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func InitializeWbsocketServer() {

	displayInterfaces()

	http.HandleFunc("/state", httpGetHandler(connect4.StateService))

	http.Handle("/echo", &WsHandler{
		ServiceFunction:      echo.EchoService,
		ErrorFunction:        echo.EchoError,
		ReconnectionFunction: echo.EchoReconnection,
	})

	http.Handle("/connect4", &WsHandler{
		ServiceFunction:      connect4.ProcessMessage,
		ErrorFunction:        connect4.ProcessError,
		ReconnectionFunction: connect4.ProcessReconnection,
	})

	err := http.ListenAndServe("0.0.0.0:1234", nil)
	utils.Assert(err)
}

func displayInterfaces() {
	fmt.Println("Starting Websocket Server")

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		addrs, _ := interf.Addrs()

		for _, addr := range addrs {

			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.To4() == nil {
				continue
			}

			fmt.Println(ip.String())
		}
	}
}
