package chat

import (
	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
)

func (c *Client) handleDispatch(message *models.WebSocketData) error {
	if c.user == nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeUnauthorized,
			Message: customErrors.ErrorMessageUnauthorized,
		}
	}

	if message.Data == nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeSyntaxError,
			Message: customErrors.GetJSONTypeIsNullErrorMessage("data"),
		}
	}

	if message.Type == "" {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeSyntaxError,
			Message: customErrors.GetJSONTypeIsNullErrorMessage("type"),
		}
	}

	switch message.Type {
	case models.EventTypeChannelJoin, models.EventTypeChannelLeave:
		return c.handleChannelEvent(message)
	default:
		return nil
	}
}
