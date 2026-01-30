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
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/room"
)

func (c *Client) handleMessageCreateEvent(message *models.WebSocketData) error {
	var data models.MessageCreateEvent

	if err := json.Unmarshal(*message.Data, &data); err != nil {
		return customErrors.GetUnmarshalError(err)
	}

	data.Author = c.user

	var roomToSend *ent.Room
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch data.Channel.RoomType {
	case room.RoomTypeDM:
		roomToSend, err = repository.GetDatabase().Rooms.GetDM(ctx, repository.GetDmID(c.user.ID, data.Channel.ID))
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

	}

	if c.currentRoomID != roomToSend.ID {
		return &customErrors.APIError{
			Code:    customErrors.ErrorCodeNotfound,
			Message: "user is not in the channel",
		}
	}

	createdMessage, err := repository.GetDatabase().Messages.Create(ctx, data)
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
