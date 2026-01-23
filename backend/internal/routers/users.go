package routers

import (
	"github.com/Migan178/misschord-backend/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func setupUsers(rg *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	public := rg.Group("/users")
	{
		public.POST("/", handler.CreateUser(authMiddleware))
		public.POST("/login", authMiddleware.LoginHandler)
	}

	private := rg.Group("/users")
	private.Use(authMiddleware.MiddlewareFunc())
	{
		private.GET("/me", handler.Me)
	}
}
