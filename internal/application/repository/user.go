package repository

import (
	"beep/internal/types"
	"context"
)

type UserRepoImpl struct{}

func (u *UserRepoImpl) Create(ctx context.Context, user *types.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepoImpl) Update(ctx context.Context, user *types.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepoImpl) Delete(ctx context.Context, user *types.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepoImpl) FindById(ctx context.Context, userId int64) (*types.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	//TODO implement me
	panic("implement me")
}
