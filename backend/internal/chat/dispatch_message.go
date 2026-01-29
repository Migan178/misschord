package chat

import (
	"context"
	"encoding/json"
	"strconv"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
)

func (c *Client) handleMessageCreateEvent(message *models.WebSocketData) error {
	var data models.MessageCreateEvent

	if err := json.Unmarshal(*message.Data, &data); err != nil {
		return customErrors.GetUnmarshalError(err)
	}

	data.Author = c.user

	switch data.Channel.ChannelType {
	case models.ChannelTypeDM:
		channelIDInt, err := strconv.Atoi(data.Channel.ID)
		if err != nil {
			return &customErrors.APIError{
				Code:    customErrors.ErrorCodeSyntaxError,
				Message: "failed to convert to int",
			}
		}
		data.Channel.InternalID = getDmID(c.user.ID, channelIDInt)
	}

	if c.currentRoomID != data.Channel.InternalID {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeNotfound,
			Message: "user is not in the channel",
		}
	}

	createdMessage, err := repository.GetDatabase().Messages.Create(context.Background(), data)
	if err != nil {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeInternalError,
			Message: customErrors.ErrorMessageInternalDBError,
		}
	}

	data.ID = &createdMessage.ID
	data.CreatedAt = &createdMessage.CreatedAt

	c.hub.parseDataAndBroadcast(message, &data)
	return nil
}
