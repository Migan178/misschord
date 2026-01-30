package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

func CreateUser(authMiddleware *jwt.GinJWTMiddleware) func(c *gin.Context) {
	return func(c *gin.Context) {
		var createData models.CreateUserRequest

		if err := c.ShouldBindJSON(&createData); err != nil {
			c.Error(err)
			return
		}

		user, err := repository.GetDatabase().Users.Create(c.Request.Context(), createData)
		if err != nil {
			c.Error(err)
			return
		}

		token, err := authMiddleware.TokenGenerator(c.Request.Context(), &models.UserToken{
			ID:    strconv.Itoa(user.ID),
			Email: user.Email,
		})
		if err != nil {
			c.Error(err)
			return
		}

		authMiddleware.SetCookie(c, token.AccessToken)
		authMiddleware.SetRefreshTokenCookie(c, token.RefreshToken)

		c.JSON(http.StatusCreated, gin.H{
			"message":      "create user and login user is success",
			"token":        token.AccessToken,
			"refreshToken": token.RefreshToken,
			"expiresAt":    token.ExpiresAt,
		})
	}
}

func CreateDM(c *gin.Context) {
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))

	var createData models.CreateDMRequest

	if err := c.ShouldBindJSON(&createData); err != nil {
		c.Error(err)
		return
	}

	room, err := repository.GetDatabase().Rooms.CreateDM(c.Request.Context(), userID, createData.RecipientID)
	if err == nil {
		c.JSON(http.StatusOK, room)
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

	c.JSON(http.StatusOK, room)
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	user, err := repository.GetDatabase().Users.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func Me(c *gin.Context) {
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))

	user, err := repository.GetDatabase().Users.Get(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {}

func DeleteUser(c *gin.Context) {}
