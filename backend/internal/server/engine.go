package server

import (
	"sync"
	"time"

	"github.com/Migan178/misschord-backend/internal/middlewares"
	"github.com/Migan178/misschord-backend/internal/routers"
	"github.com/gin-gonic/gin"
)

var instance *gin.Engine
var once sync.Once

func GetEngine() *gin.Engine {
	once.Do(func() {
		instance = gin.Default()

		instance.Use(middlewares.TimeoutMiddleWare(time.Second * 5))

		routers.SetupRouter(instance)
	})

	return instance
}
