package repository

import (
	"context"
	"errors"
	"fmt"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/room"
	"github.com/go-sql-driver/mysql"
)

func GetDmID(id1, id2 int) string {
	if id1 < id2 {
		return fmt.Sprintf("%d:%d", id1, id2)
	}

	return fmt.Sprintf("%d:%d", id2, id1)
}

type RoomRepository struct {
	client *ent.Client
}

func newRoomRepository(client *ent.Client) *RoomRepository {
	return &RoomRepository{client}
}

func (r *RoomRepository) CreateDM(ctx context.Context, userID, recipientID int) (*ent.Room, error) {
	room, err := r.client.Room.Create().
		SetRoomType(room.RoomTypeDM).
		SetDmKey(GetDmID(userID, recipientID)).
		AddMemberIDs(userID, recipientID).
		Save(ctx)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return nil, customErrors.ErrDuplicatedUniqueValue
			}
		}
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) GetDM(ctx context.Context, dmKey string) (*ent.Room, error) {
	channel, err := r.client.Room.Query().
		Where(room.DmKey(dmKey)).
		Only(ctx)
	if errors.As(err, new(*ent.NotFoundError)) {
		return nil, customErrors.ErrNoUser
	}

	return channel, nil
}
