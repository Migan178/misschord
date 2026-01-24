package routers

import (
	"github.com/Migan178/misschord-backend/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := app.Group("/v1")
	{
		v1.GET("/ws", handler.WSRoot)
		setupUsers(v1, authMiddleware)
	}
}
