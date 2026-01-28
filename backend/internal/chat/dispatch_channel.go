package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
)

func getDmID(id1, id2 int) string {
	if id1 < id2 {
		return fmt.Sprintf("%d:%d", id1, id2)
	}

	return fmt.Sprintf("%d:%d", id2, id1)
}

func (c *Client) handleChannelEvent(message *models.WebSocketData) error {
	var data models.ChannelData

	if err := json.Unmarshal(*message.Data, &data); err != nil {
		return customErrors.GetUnmarshalError(err)
	}

	if data.ID == fmt.Sprintf("%d", c.user.ID) {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeInvalidValue,
			Message: customErrors.ErrorMessageChannelIDIsInvalid,
		}
	}

	if data.ChannelType == "" {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeSyntaxError,
			Message: customErrors.GetJSONTypeIsNullErrorMessage("type"),
		}
	}

	channelIDInt, err := strconv.Atoi(data.ID)
	if err != nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeSyntaxError,
			Message: "failed to convert to int",
		}
	}

	var roomID string

	switch data.ChannelType {
	case models.ChannelTypeDM:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		targetUser, err := repository.GetDatabase().Users.Get(ctx, channelIDInt)
		if err != nil {
			if errors.As(err, new(*ent.NotFoundError)) {
				return &customErrors.APIError{
					Code:    customErrors.ErrorCodeNotfound,
					Message: fmt.Sprintf("user %d is not found", channelIDInt),
				}
			}

			return &customErrors.APIError{
				Code:    customErrors.ErrorCodeInternalError,
				Message: customErrors.ErrorMessageInternalDBError,
			}
		}

		if c.user.ID < targetUser.ID {
			roomID = fmt.Sprintf("%d:%d", c.user.ID, targetUser.ID)
		} else {
			roomID = fmt.Sprintf("%d:%d", targetUser.ID, c.user.ID)
		}
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
