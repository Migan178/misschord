package routers

import (
	"github.com/Migan178/misschord-backend/internal/handler"
	"github.com/gin-gonic/gin"
)

func setupUsers(rg *gin.RouterGroup) {
	public := rg.Group("/users")
	{
		public.POST("/", handler.CreateUser)
	}
}
