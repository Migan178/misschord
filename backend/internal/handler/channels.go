package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func CreateDM(c *gin.Context) {
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))

	var createData models.CreateDMRequest

	if err := c.ShouldBindJSON(&createData); err != nil {
		c.Error(err)
		return
	}

	room, err := repository.GetDatabase().Rooms.CreateDM(c.Request.Context(), userID, createData.RecipientID)
	if err == nil {
		c.JSON(http.StatusCreated, room)
		return
	}

	var dbErr *repository.DatabaseError
	if !errors.As(err, &dbErr) {
		c.Error(err)
		return
	}

	if dbErr.Code != repository.ErrorCodeNotFound {
		c.Error(err)
		return
	}

	dmKey := repository.GetDmID(userID, createData.RecipientID)
	room, err = repository.GetDatabase().Rooms.GetDM(c.Request.Context(), dmKey)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, room)
}

func CreateChannel(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "wip"})
}

func GetChannel(c *gin.Context) {
	path := c.Request.URL.Path
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))
	channelID, err := strconv.Atoi(c.Param("channelID"))
	if err != nil {
		c.Error(err)
		return
	}

	var room *ent.Room

	if strings.Contains(path, "/users/me") {
		room, err = repository.GetDatabase().Rooms.GetDM(c.Request.Context(), repository.GetDmID(userID, channelID))
	} else {
		room, err = repository.GetDatabase().Rooms.GetRoom(c.Request.Context(), channelID)
	}

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, room)
}
