package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func HandleGetMessages(c *gin.Context) {
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))
	path := c.Request.URL.Path
	channelID, err := strconv.Atoi(c.Param("channelID"))
	if err != nil {
		c.Error(err)
		return
	}

	var messages []*models.MessageCreateEvent

	if strings.Contains(path, "/users/me") {
		messages, err = repository.GetDatabase().Messages.GetDmMessages(c.Request.Context(), repository.GetDmID(userID, channelID))
	} else {

	}

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, messages)
}
