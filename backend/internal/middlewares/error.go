package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			msgs := make(map[string]string)

			for _, err := range validateErr {
				msgs[err.Field()] = customErrors.GetErrorMessage(err)
			}

			c.JSON(http.StatusBadRequest, customErrors.APIError{
				Code:     customErrors.ErrorCodeSyntaxError,
				Messages: msgs,
			})
			return
		}

		if errors.As(err, new(*json.UnmarshalTypeError)) || errors.As(err, new(*json.SyntaxError)) {
			c.JSON(http.StatusBadRequest, customErrors.GetUnmarshalError(err))
			return
		}

		var dbErr *repository.DatabaseError
		if errors.As(err, &dbErr) {
			switch dbErr.Code {
			case repository.ErrorCodeAuthenticationFailed:
				c.JSON(http.StatusUnauthorized, customErrors.APIError{
					Code:    customErrors.ErrorCodeAuthenticationErr,
					Message: "user is not valid",
				})
				return
			case repository.ErrorCodeNotFound:
				c.JSON(http.StatusBadRequest, customErrors.APIError{
					Code:    customErrors.ErrorCodeNotfound,
					Message: "not found",
				})
				return
			case repository.ErrorCodeConstraint:
				fmt.Printf("constraint error: %v\n", err)
				c.JSON(http.StatusBadRequest, customErrors.APIError{
					Code:    customErrors.ErrorCodeInvalidValue,
					Message: customErrors.ErrorMessageConstraintErr,
				})
				return
			case repository.ErrorCodeOther:
				fmt.Printf("db err: %v\n", err)
				c.JSON(http.StatusInternalServerError, customErrors.APIError{
					Code:    customErrors.ErrorCodeInternalError,
					Message: customErrors.ErrorMessageInternalDBError,
				})
				return
			}
		}

		if errors.Is(err, strconv.ErrSyntax) || errors.Is(err, strconv.ErrRange) {
			c.JSON(http.StatusBadRequest, customErrors.APIError{
				Code:    customErrors.ErrorCodeSyntaxError,
				Message: err.Error(),
			})
			return
		}

		fmt.Printf("unexpected err: %v\n", err)
		c.JSON(http.StatusInternalServerError, customErrors.APIError{
			Code:    customErrors.ErrorCodeInternalError,
			Message: "an error occurred",
		})
	}
}
