package repository

import (
	"context"

	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/message"
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
		SetChannelID(data.Channel.InternalID).
		SetChannelType(message.ChannelType(data.Channel.ChannelType)).
		Save(ctx)
}

func (r *MessageRepository) GetMessagesByChannel(ctx context.Context, data models.ChannelData) ([]*models.MessageCreateEvent, error) {
	messages, err := r.client.Message.Query().
		Where(
			message.ChannelID(data.InternalID),
			message.ChannelTypeEQ(message.ChannelType(data.ChannelType)),
		).
		WithAuthor().
		All(ctx)
	if err != nil {
		return make([]*models.MessageCreateEvent, 0), err
	}

	var messagesToReturn []*models.MessageCreateEvent

	for _, message := range messages {
		messagesToReturn = append(messagesToReturn, returnToMessageCreateEvent(message, data))
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
