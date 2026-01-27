package repository

import (
	"context"
	"errors"

	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/user"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	client *ent.Client
}

func newUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client}
}

func (r *UserRepository) Create(ctx context.Context, data models.CreateUserRequest) (*ent.User, error) {
	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user, err := r.client.User.Create().
		SetHandle(data.Handle).
		SetEmail(data.Email).
		SetHashedPassword(hashedPassword).
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

	return user, nil
}

func (r *UserRepository) Get(ctx context.Context, id int) (*ent.User, error) {
	user, err := r.client.User.Get(ctx, id)
	if err != nil {
		if _, ok := err.(*ent.NotFoundError); ok {
			return nil, customErrors.ErrNoUser
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	user, err := r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)
	if err != nil {
		if _, ok := err.(*ent.NotFoundError); ok {
			return nil, jwt.ErrFailedAuthentication
		}

		return nil, err
	}

	return user, nil
}
