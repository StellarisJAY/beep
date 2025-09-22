package interfaces

import (
	"beep/internal/types"
	"context"
)

type UserRepo interface {
	Create(ctx context.Context, user *types.User) error
	Update(ctx context.Context, user *types.User) error
	Delete(ctx context.Context, user *types.User) error
	FindById(ctx context.Context, userId int64) (*types.User, error)
	FindByEmail(ctx context.Context, email string) (*types.User, error)
}

type UserService interface {
	// TODO user service
}
