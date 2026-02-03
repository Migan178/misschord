package routers

import (
	"github.com/Migan178/misschord-backend/internal/chat"
	"github.com/Migan178/misschord-backend/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *gin.Engine, authMiddleware *jwt.GinJWTMiddleware, hub *chat.Hub) {
	v1 := app.Group("/v1")
	{
		v1.GET("/ws", func(c *gin.Context) {
			handler.ServeWS(hub, authMiddleware, c.Writer, c.Request)
		})

		setupUsers(v1, authMiddleware, hub)
	}
}
