package main

import (
	"fmt"

	"github.com/Migan178/misschord-backend/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	app.GET("/", handler.Root)
	app.GET("/ws", handler.WSRoot)

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
