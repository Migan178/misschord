package routers

import (
	"strings"

	"github.com/Migan178/misschord-backend/internal/chat"
	"github.com/Migan178/misschord-backend/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func setChannels(rg *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware, hub *chat.Hub) {
	channels := rg.Group("/channels")
	channels.Use(authMiddleware.MiddlewareFunc())
	{
		channels.POST("", func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.Contains(path, "/users/me") {
				handler.HandleCreateDM(c)
			} else {
				handler.HandleCreateChannel(c)
			}
		})
		channels.GET("/:channelID", handler.HandleGetChannel)
	}

	setupMessages(rg.Group("/channels/:channelID"), authMiddleware, hub)
}
