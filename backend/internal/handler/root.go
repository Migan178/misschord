package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"Hello": "World!"})
}
