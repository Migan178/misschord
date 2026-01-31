package routers

import (
	"strings"

	"github.com/Migan178/misschord-backend/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func setChannels(rg *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	channels := rg.Group("/channels")
	channels.Use(authMiddleware.MiddlewareFunc())
	{
		channels.POST("/", func(c *gin.Context) {
			if strings.Contains(c.Request.URL.Path, "/users/me") {
				handler.CreateDM(c)
			} else {
				handler.CreateChannel(c)
			}
		})
		channels.GET("/:id", handler.GetChannel)
	}
}
