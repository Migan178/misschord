package repository

import (
	"context"

	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	"github.com/Migan178/misschord-backend/internal/repository/ent/user"
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
		code := ErrorCodeOther

		if ent.IsConstraintError(err) {
			code = ErrorCodeConstraint
		}

		return nil, &DatabaseError{
			Code:   code,
			RawErr: err,
		}
	}

	return user, nil
}

func (r *UserRepository) Get(ctx context.Context, id int) (*ent.User, error) {
	user, err := r.client.User.Get(ctx, id)
	if err != nil {
		code := ErrorCodeOther

		if ent.IsNotFound(err) {
			code = ErrorCodeAuthenticationFailed
		}

		return nil, &DatabaseError{
			Code:   code,
			RawErr: err,
		}
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	user, err := r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)
	if err != nil {
		code := ErrorCodeOther

		if ent.IsNotFound(err) {
			code = ErrorCodeNotFound
		}

		return nil, &DatabaseError{
			Code:   code,
			RawErr: err,
		}
	}

	return user, nil
}
