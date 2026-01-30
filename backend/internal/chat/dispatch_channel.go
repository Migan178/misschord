package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	"github.com/Migan178/misschord-backend/internal/repository/ent/room"
)

func (c *Client) handleChannelEvent(message *models.WebSocketData) error {
	var data models.ChannelData

	if err := json.Unmarshal(*message.Data, &data); err != nil {
		return customErrors.GetUnmarshalError(err)
	}

	if data.ID == c.user.ID {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeInvalidValue,
			Message: customErrors.ErrorMessageChannelIDIsInvalid,
		}
	}

	if data.RoomType == "" {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeSyntaxError,
			Message: customErrors.GetJSONTypeIsNullErrorMessage("type"),
		}
	}

	var roomID int

	switch data.RoomType {
	case room.RoomTypeDM:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		room, err := repository.GetDatabase().Rooms.GetDM(ctx, repository.GetDmID(c.user.ID, data.ID))
		if err != nil {
			if errors.Is(err, customErrors.ErrNoUser) {
				return &customErrors.APIError{
					Code:    customErrors.ErrorCodeNotfound,
					Message: fmt.Sprintf("user %d is not found", data.ID),
				}
			}

			return &customErrors.APIError{
				Code:    customErrors.ErrorCodeInternalError,
				Message: customErrors.ErrorMessageInternalDBError,
			}
		}

		roomID = room.ID
	}

	switch message.Type {
	case models.EventTypeChannelJoin:
		select {
		case c.hub.join <- joinOrLeaveData{c, roomID}:
			c.parseDataAndSend(message, data)
		case <-time.After(writeWait + 5*time.Second):
			return customErrors.ErrFailedToSend
		}
	case models.EventTypeChannelLeave:
		select {
		case c.hub.leave <- joinOrLeaveData{c, roomID}:
			c.parseDataAndSend(message, data)
		case <-time.After(writeWait + 5*time.Second):
			return customErrors.ErrFailedToSend
		}
	}

	return nil
}
