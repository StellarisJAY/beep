package types

import "time"

// Workspace 工作空间
type Workspace struct {
	BaseEntity
	Name        string `json:"name" gorm:"not null;type:varchar(64);'"`
	Description string `json:"description" gorm:"not null;type:varchar(255)"`
}

func (Workspace) TableName() string {
	return "workspaces"
}

type WorkspaceRole string

const (
	WorkspaceRoleOwner  WorkspaceRole = "owner"  // 工作空间拥有者，最高权限
	WorkspaceRoleAdmin  WorkspaceRole = "admin"  // 工作空间管理员，最高权限
	WorkspaceRoleEditor WorkspaceRole = "editor" // 工作空间编辑者，可以编辑工作空间下的所有智能体和知识库
	WorkspaceRoleNormal WorkspaceRole = "normal" // 普通权限，只能使用智能体和读取知识库
)

// UserWorkspace 用户-工作空间关联
type UserWorkspace struct {
	BaseEntity
	UserId      int64         `json:"user_id" gorm:"not null;"`
	WorkspaceId int64         `json:"workspace_id" gorm:"not null;"`
	Role        WorkspaceRole `json:"role" gorm:"default:'owner';not null;type:varchar(16)"`
}

// WorkspaceMember 工作空间成员
type WorkspaceMember struct {
	Id            int64         `json:"id,string"`
	Name          string        `json:"name"`
	Email         string        `json:"email"`
	Role          WorkspaceRole `json:"role"`
	LastLoginTime *time.Time    `json:"last_login_time"`
}

// InviteWorkspaceMemberReq 邀请工作空间成员
type InviteWorkspaceMemberReq struct {
	WorkspaceId int64         `json:"workspace_id,string"`
	Emails      []string      `json:"emails"` // 邀请的邮箱列表
	Role        WorkspaceRole `json:"role"`   // 加入后的角色
}

type SetWorkspaceRoleReq struct {
	WorkspaceId int64         `json:"workspace_id,string"`
	UserId      int64         `json:"user_id,string"`
	Role        WorkspaceRole `json:"role"`
}
