package routers

import (
	"github.com/Migan178/misschord-backend/internal/chat"
	"github.com/Migan178/misschord-backend/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func setupMessages(rg *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware, hub *chat.Hub) {
	messages := rg.Group("/messages")
	messages.Use(authMiddleware.MiddlewareFunc())

	messageHandler := &handler.MessageHandler{
		Hub: hub,
	}

	{
		messages.POST("", messageHandler.HandleCreateMessage)
		messages.GET("", handler.HandleGetMessages)
	}
}
