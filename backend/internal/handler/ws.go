package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSRoot(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}

	defer ws.Close()

	if err = ws.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!")); err != nil {
		fmt.Println("error in writing data:", err)
		return
	}

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("error in reading data:", err)
			break
		}

		fmt.Println("received msg:", string(msg))

		if err := ws.WriteMessage(websocket.TextMessage, []byte("received: "+string(msg))); err != nil {
			fmt.Println("error in writing data:", err)
			break
		}
	}
}
