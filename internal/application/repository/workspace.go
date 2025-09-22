package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type WorkspaceRepo struct {
	db *gorm.DB
}

func NewWorkspaceRepo(db *gorm.DB) interfaces.WorkspaceRepo {
	return &WorkspaceRepo{
		db: db,
	}
}

func (w *WorkspaceRepo) Create(ctx context.Context, workspace *types.Workspace) error {
	return w.db.WithContext(ctx).Create(workspace).Error
}

func (w *WorkspaceRepo) Update(ctx context.Context, workspace *types.Workspace) error {
	//TODO implement me
	panic("implement me")
}

func (w *WorkspaceRepo) FindById(ctx context.Context, id int64) (*types.Workspace, error) {
	//TODO implement me
	panic("implement me")
}
