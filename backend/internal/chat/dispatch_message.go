package chat

import (
	"encoding/json"
	"strconv"
	"time"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
)

func (c *Client) handleMessageCreateEvent(message *models.WebSocketData) error {
	var data models.MessageCreateEvent

	if err := json.Unmarshal(*message.Data, &data); err != nil {
		return customErrors.GetUnmarshalError(err)
	}

	now := time.Now()
	data.CreatedAt = &now
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

	c.hub.parseDataAndBroadcast(message, &data)
	return nil
}
