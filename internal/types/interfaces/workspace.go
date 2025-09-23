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
	ListMember(ctx context.Context, workspaceId int64) ([]*types.WorkspaceMember, error)
	FindUserDefaultWorkspace(ctx context.Context, userId int64) (*types.Workspace, error)
	FindUserJoinedWorkspace(ctx context.Context, userId int64) ([]*types.Workspace, error)
	Update(ctx context.Context, uw *types.UserWorkspace) error
}

// WorkspaceService 工作空间服务
type WorkspaceService interface {
	// FindById 根据ID查询工作空间
	FindById(ctx context.Context, id int64) (*types.Workspace, error)
	// ListMembers 查询工作空间成员
	ListMembers(ctx context.Context, workspaceId int64) ([]*types.WorkspaceMember, error)
	// ListByUserId 查询用户加入的工作空间
	ListByUserId(ctx context.Context, userId int64) ([]*types.Workspace, error)
	// InviteMember 邀请工作空间成员
	InviteMember(ctx context.Context, req types.InviteWorkspaceMemberReq) error
	// SwitchWorkspace 切换工作空间
	SwitchWorkspace(ctx context.Context, id int64) error
	// SetRole 设置工作空间成员角色
	SetRole(ctx context.Context, req types.SetWorkspaceRoleReq) error
}
