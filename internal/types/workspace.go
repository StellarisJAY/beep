package types

// Workspace 工作空间
type Workspace struct {
	BaseEntity
	Name        string `json:"name"`
	Description string `json:"description"`
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
	UserId      uint64        `json:"user_id"`
	WorkspaceId uint64        `json:"workspace_id"`
	Role        WorkspaceRole `json:"role"`
}
