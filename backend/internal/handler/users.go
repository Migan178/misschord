package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUser(authMiddleware *jwt.GinJWTMiddleware) func(c *gin.Context) {
	return func(c *gin.Context) {
		var createData models.CreateUserRequest

		if err := c.ShouldBindJSON(&createData); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				msgs := make(map[string]string)

				for _, err := range errs {
					msgs[err.Field()] = errors.GetErrorMessage(err)
				}

				c.JSON(http.StatusBadRequest, gin.H{"errors": msgs})
				return
			}

			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "an error occurred"})
			return
		}

		user, err := repository.GetDatabase().Users.Create(c.Request.Context(), createData)
		if err != nil {
			if err == errors.ErrDuplicatedUniqueValue {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred"})
			return
		}

		token, err := authMiddleware.TokenGenerator(c.Request.Context(), &models.UserToken{
			ID:    strconv.Itoa(user.ID),
			Email: user.Email,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred in generating token. but user is created."})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":       "create user and login user is success",
			"token":         token.AccessToken,
			"refresh_token": token.RefreshToken,
			"expire":        token.ExpiresAt,
		})
	}
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is not valid",
		})
		return
	}

	user, err := repository.GetDatabase().Users.Get(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrNoUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return

		}

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "an error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Me(c *gin.Context) {
	userID, _ := strconv.Atoi(jwt.ExtractClaims(c)["id"].(string))

	user, err := repository.GetDatabase().Users.Get(c.Request.Context(), userID)
	if err != nil {
		if err == errors.ErrNoUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "an error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {}

func DeleteUser(c *gin.Context) {}
