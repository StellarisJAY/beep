package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type UserWorkspaceRepo struct {
	db *gorm.DB
}

func NewUserWorkspaceRepo(db *gorm.DB) interfaces.UserWorkspaceRepo {
	return &UserWorkspaceRepo{
		db: db,
	}
}

func (u *UserWorkspaceRepo) Create(ctx context.Context, uw *types.UserWorkspace) error {
	return u.db.WithContext(ctx).Create(uw).Error
}

func (u *UserWorkspaceRepo) FindByUser(ctx context.Context, userId int64) ([]*types.UserWorkspace, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserWorkspaceRepo) Find(ctx context.Context, workspaceId, userId int64) (*types.UserWorkspace, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserWorkspaceRepo) FindByWorkspace(ctx context.Context, workspaceId int64) (*types.UserWorkspace, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserWorkspaceRepo) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
