package repository

import (
	"context"

	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/room"
)

type MessageRepository struct {
	client *ent.Client
}

func newMessageRepository(client *ent.Client) *MessageRepository {
	return &MessageRepository{client}
}

func (r *MessageRepository) Create(ctx context.Context, data models.MessageCreateEvent) (*ent.Message, error) {
	return r.client.Message.Create().
		SetAuthorID(data.Author.ID).
		SetMessage(data.Message).
		SetRoomID(data.Channel.ID).
		Save(ctx)
}

func (r *MessageRepository) GetDmMessages(ctx context.Context, dmKey string) ([]*models.MessageCreateEvent, error) {
	messages, err := r.client.Room.Query().
		Where(room.DmKey(dmKey)).
		QueryMessages().
		All(ctx)
	if err != nil {
		return make([]*models.MessageCreateEvent, 0), err
	}

	var messagesToReturn []*models.MessageCreateEvent

	for _, message := range messages {
		messagesToReturn = append(messagesToReturn, returnToMessageCreateEvent(message, models.ChannelData{
			ID:       message.RoomID,
			RoomType: room.RoomTypeDM,
		}))
	}

	return messagesToReturn, nil
}

func returnToMessageCreateEvent(message *ent.Message, channel models.ChannelData) *models.MessageCreateEvent {
	return &models.MessageCreateEvent{
		ID:        &message.ID,
		Author:    message.Edges.Author,
		Message:   message.Message,
		Channel:   channel,
		CreatedAt: &message.CreatedAt,
	}
}
