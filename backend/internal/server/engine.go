package server

import (
	"sync"
	"time"

	"github.com/Migan178/misschord-backend/internal/middlewares"
	"github.com/Migan178/misschord-backend/internal/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var instance *gin.Engine
var once sync.Once

func GetEngine() *gin.Engine {
	once.Do(func() {
		instance = gin.Default()

		authMiddleware, err := middlewares.NewAuthMiddleware()
		if err != nil {
			panic(err)
		}

		if err = authMiddleware.MiddlewareInit(); err != nil {
			panic(err)
		}

		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"http://localhost:3000"}
		config.AllowCredentials = true

		instance.Use(cors.New(config))
		instance.Use(middlewares.TimeoutMiddleWare(time.Second * 5))

		routers.SetupRouter(instance, authMiddleware)
	})

	return instance
}
