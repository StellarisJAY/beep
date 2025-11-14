package interfaces

import (
	"beep/internal/types"
	"context"
)

// WorkspaceRepo 工作空间数据库
type WorkspaceRepo interface {
	// Create 创建工作空间
	Create(ctx context.Context, workspace *types.Workspace) error
	// Update 更新工作空间
	Update(ctx context.Context, workspace *types.Workspace) error
	// FindById 根据ID查询工作空间
	FindById(ctx context.Context, id string) (*types.Workspace, error)
}

// UserWorkspaceRepo 用户工作空间数据库
type UserWorkspaceRepo interface {
	// Create 创建用户工作空间
	Create(ctx context.Context, uw *types.UserWorkspace) error
	// FindByUser 根据用户ID查询用户工作空间
	FindByUser(ctx context.Context, userId string) ([]*types.UserWorkspace, error)
	// Find 根据工作空间ID和用户ID查询用户工作空间
	Find(ctx context.Context, workspaceId, userId string) (*types.UserWorkspace, error)
	// FindByWorkspace 根据工作空间ID查询用户工作空间
	FindByWorkspace(ctx context.Context, workspaceId string) (*types.UserWorkspace, error)
	// Delete 删除用户工作空间
	Delete(ctx context.Context, id string) error
	// ListMember 查询工作空间成员
	ListMember(ctx context.Context, workspaceId string) ([]*types.WorkspaceMember, error)
	// FindUserDefaultWorkspace 根据用户ID查询用户默认工作空间
	FindUserDefaultWorkspace(ctx context.Context, userId string) (*types.Workspace, error)
	// FindUserJoinedWorkspace 根据用户ID查询用户加入的工作空间
	FindUserJoinedWorkspace(ctx context.Context, userId string) ([]*types.Workspace, error)
	// Update 更新用户工作空间
	Update(ctx context.Context, uw *types.UserWorkspace) error
}

// WorkspaceService 工作空间服务
type WorkspaceService interface {
	// FindById 根据ID查询工作空间
	FindById(ctx context.Context, id string) (*types.Workspace, error)
	// ListMembers 查询工作空间成员
	ListMembers(ctx context.Context, workspaceId string) ([]*types.WorkspaceMember, error)
	// ListByUserId 查询用户加入的工作空间
	ListByUserId(ctx context.Context, userId string) ([]*types.Workspace, error)
	// InviteMember 邀请工作空间成员
	InviteMember(ctx context.Context, req types.InviteWorkspaceMemberReq) error
	// SwitchWorkspace 切换工作空间
	SwitchWorkspace(ctx context.Context, id string) error
	// SetRole 设置工作空间成员角色
	SetRole(ctx context.Context, req types.SetWorkspaceRoleReq) error
}
