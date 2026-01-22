package repository

import (
	"context"
	"errors"
	"time"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	client *ent.Client
}

func newUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client}
}

type CreateUserRequest struct {
	Handle         string `json:"handle" binding:"required,min=4,max=16"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required"`
	PasswordCheck  string `json:"password_check" binding:"required,eqfield=Password"`
	hashedPassword string `json:"-"`
}

type User struct {
	ID      int    `json:"id"`
	Profile string `json:"profile"`

	// Unique
	Handle string `json:"handle"`
	// Unique
	Email string `json:"-"`

	HashedPassword string    `json:"-"`
	Description    *string   `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

func (c *UserRepository) Create(ctx context.Context, data CreateUserRequest) (*ent.User, error) {
	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user, err := c.client.User.Create().
		SetHandle(data.Handle).
		SetEmail(data.Email).
		SetHashedPassword(hashedPassword).
		Save(ctx)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return nil, customErrors.DuplicatedUniqueValueErr
			}
		}
		return nil, err
	}

	return user, nil
}
