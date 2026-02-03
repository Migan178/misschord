package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Migan178/misschord-backend/internal/chat"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	"github.com/Migan178/misschord-backend/internal/repository/ent/room"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	Hub *chat.Hub
}

func (h *MessageHandler) HandleCreateMessage(c *gin.Context) {
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))
	path := c.Request.URL.Path

	channelID, err := strconv.Atoi(c.Param("channelID"))
	if err != nil {
		c.Error(err)
		return
	}

	var createData models.MessageCreateData

	if err := c.ShouldBindBodyWithJSON(&createData); err != nil {
		c.Error(err)
		return
	}

	channel := &models.ChannelData{
		ID: channelID,
	}

	if strings.Contains(path, "/users/me") {
		channel.RoomType = room.RoomTypeDM
	} else {
		channel.RoomType = room.RoomTypeCHANNEL
	}

	createData.Channel = channel
	message, err := repository.GetDatabase().Messages.Create(c.Request.Context(), userID, createData)
	if err != nil {
		c.Error(err)
		return
	}

	op := models.OPCodeDispatch
	h.Hub.ParseDataAndBroadcast(&models.WebSocketData{
		OP:   &op,
		Type: models.EventTypeMessageCreate,
	}, message)

	c.JSON(http.StatusCreated, message)
}

func HandleGetMessages(c *gin.Context) {
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))
	path := c.Request.URL.Path
	channelID, err := strconv.Atoi(c.Param("channelID"))
	if err != nil {
		c.Error(err)
		return
	}

	var messages []*models.MessageResponse

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
