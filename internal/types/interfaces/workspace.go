package interfaces

import (
	"beep/internal/types"
	"context"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, workspace *types.Workspace) error
	Update(ctx context.Context, workspace *types.Workspace) error
	FindById(ctx context.Context, id int64) (*types.Workspace, error)
}

type UserWorkspaceRepo interface {
	Create(ctx context.Context, uw *types.UserWorkspace) error
	FindByUser(ctx context.Context, userId int64) ([]*types.UserWorkspace, error)
	Find(ctx context.Context, workspaceId, userId int64) (*types.UserWorkspace, error)
	FindByWorkspace(ctx context.Context, workspaceId int64) (*types.UserWorkspace, error)
	Delete(ctx context.Context, id int64) error
}
