package main

import (
	"fmt"

	"github.com/Migan178/misschord-backend/internal/configs"
	"github.com/Migan178/misschord-backend/internal/server"
)

func main() {
	app := server.GetEngine()

	if err := app.Run(fmt.Sprintf(":%d", configs.GetConfig().Backend.Port)); err != nil {
		fmt.Println(err)
	}
}
