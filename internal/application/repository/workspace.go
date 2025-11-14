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
	return w.db.WithContext(ctx).Updates(workspace).Error
}

func (w *WorkspaceRepo) FindById(ctx context.Context, id string) (*types.Workspace, error) {
	var workspace *types.Workspace
	err := w.db.WithContext(ctx).Model(&types.Workspace{}).
		Where("id = ?", id).
		First(&workspace).Error
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

// workspaceScope gorm 工作空间查询注入
func workspaceScope(ctx context.Context) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if workspaceId, ok := ctx.Value(types.WorkspaceIdContextKey).(string); ok {
			db = db.Where("workspace_id = ?", workspaceId)
		}
		return db
	}
}

func workspaceScopeWithTable(ctx context.Context, table string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if workspaceId, ok := ctx.Value(types.WorkspaceIdContextKey).(string); ok {
			db = db.Where(table+".workspace_id = ?", workspaceId)
		}
		return db
	}
}
