package chat

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	jwt "github.com/appleboy/gin-jwt/v3"
)

func (c *Client) handleIdentify(message *models.WebSocketData) error {
	if c.user != nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeAuthorizationError,
			Message: customErrors.ErrorMessageAlreadyIdentified,
		}
	}

	if message.Data == nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeSyntaxError,
			Message: customErrors.GetJSONTypeIsNullErrorMessage("data"),
		}
	}

	var data models.IdentifyData

	if err := json.Unmarshal(*message.Data, &data); err != nil {
		return customErrors.GetUnmarshalError(err)
	}

	if data.Token == "" {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeAuthorizationError,
			Message: customErrors.ErrorMessageInvalidToken,
		}
	}

	token, err := c.hub.authMiddleware.ParseTokenString(data.Token)
	if err != nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeAuthorizationError,
			Message: customErrors.ErrorMessageInvalidToken,
		}
	}

	claims := jwt.ExtractClaimsFromToken(token)
	id, err := strconv.Atoi(claims["id"].(string))
	if err != nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeAuthorizationError,
			Message: customErrors.ErrorMessageInvalidToken,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbUser, err := repository.GetDatabase().Users.Get(ctx, id)
	if err != nil {
		var dbErr *repository.DatabaseError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.ErrorCodeNotFound {
				return &customErrors.APIError{
					Code:    customErrors.ErrorCodeAuthorizationError,
					Message: "user is not found",
				}
			}
		}

		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeInternalError,
			Message: customErrors.ErrorMessageInternalDBError,
		}
	}

	c.user = dbUser
	op := models.OPCodeReady
	c.parseDataAndSend(&models.WebSocketData{OP: &op}, models.ReadyData{User: c.user})

	return nil
}
