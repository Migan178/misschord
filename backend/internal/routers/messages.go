package routers

import (
	"github.com/Migan178/misschord-backend/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func setupMessages(rg *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	messages := rg.Group("/messages")
	messages.Use(authMiddleware.MiddlewareFunc())
	{
		messages.GET("/", handler.HandleGetMessages)
	}
}
