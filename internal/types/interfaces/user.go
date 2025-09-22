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
	CheckPassword(ctx context.Context, email string, password string) (*types.User, error)
}

type UserService interface {
	Register(ctx context.Context, req types.RegisterReq) error
	Login(ctx context.Context, req types.LoginReq) (*types.User, error)
}
