package routers

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	v1 := app.Group("/v1")

	setupUsers(v1, authMiddleware)
}
