package chat

import (
	"errors"
	"io"
	"time"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
	}()

	c.conn.SetReadLimit(messageLimit)
	c.conn.SetReadDeadline(time.Now().Add(clientWait + 10*time.Second))

	for {
		var messageData models.WebSocketData

		if err := c.conn.ReadJSON(&messageData); err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				err = &customErrors.APIError{
					Code:    customErrors.ErrorCodeSyntaxError,
					Message: err.Error(),
				}
			} else {
				err = customErrors.GetUnmarshalError(err)
			}

			op := models.OPCodeError
			c.parseDataAndSend(&models.WebSocketData{OP: &op}, err)
			return
		}

		if messageData.OP == nil {
			continue
		}

		if c.user != nil {
			c.conn.SetReadDeadline(time.Now().Add(clientWait + 10*time.Second))
		}

		switch *messageData.OP {
		case models.OPCodeHeartBeat:
			if c.user == nil {
				continue
			}

			op := models.OPCodeHeartBeatACK
			c.safeSend(&models.WebSocketData{OP: &op})
		case models.OPCodeDispatch:
			if err := c.handleDispatch(&messageData); err != nil {
				if errors.Is(err, customErrors.ErrFailedToSend) {
					return
				}

				var apiErr *customErrors.APIError
				if errors.As(err, &apiErr) {
					op := models.OPCodeError
					c.parseDataAndSend(&models.WebSocketData{OP: &op}, err)

					isUnauthorized := apiErr.Code == customErrors.ErrorCodeUnauthorized
					isInternalErr := apiErr.Code == customErrors.ErrorCodeInternalError
					if isUnauthorized || isInternalErr {
						return
					}

					continue
				}
			}
		case models.OPCodeIdentify:
			if err := c.handleIdentify(&messageData); err != nil {
				op := models.OPCodeError
				messageToSend := models.WebSocketData{OP: &op}
				var apiErr *customErrors.APIError
				if errors.As(err, &apiErr) {
					c.parseDataAndSend(&messageToSend, err)

					isAuthorizationErr := apiErr.Code == customErrors.ErrorCodeAuthorizationError
					isInternalErr := apiErr.Code == customErrors.ErrorCodeInternalError
					if isAuthorizationErr || isInternalErr {
						return
					}

					continue
				}
			}
		default:
			continue
		}
	}
}
