package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"Hello": "World!"})
}
