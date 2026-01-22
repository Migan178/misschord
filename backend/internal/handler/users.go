package handler

import (
	"fmt"
	"net/http"

	"github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUser(c *gin.Context) {
	var createData repository.CreateUserRequest

	if err := c.ShouldBind(&createData); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			msgs := make(map[string]string)

			for _, err := range errs {
				msgs[err.Field()] = GetErrorMessage(err)
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
		if err == errors.DuplicatedUniqueValueErr {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func GetUser(c *gin.Context) {}

func UpdateUser(c *gin.Context) {}

func DeleteUser(c *gin.Context) {}
