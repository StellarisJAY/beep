package types

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

// MCPServer 注册MCP服务
type MCPServer struct {
	BaseEntity
	Name        string `json:"name" gorm:"type:varchar(255);not null;"`
	Url         string `json:"url" gorm:"type:varchar(255);not null;"`
	Description string `json:"description" gorm:"type:varchar(255);not null;"`
	CreateBy    int64  `json:"create_by" gorm:"type:bigint;not null;"`
	WorkspaceId int64  `json:"workspace_id" gorm:"type:bigint;not null;"`

	Available bool        `json:"available" gorm:"-"`
	Tools     []*mcp.Tool `json:"tools" gorm:"-"`
}

func (*MCPServer) TableName() string {
	return "mcp_servers"
}

func (m *MCPServer) BeforeSave(tx *gorm.DB) error {
	if err := m.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if m.WorkspaceId == 0 {
		if workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(int64); ok {
			m.WorkspaceId = workspaceId
		}
	}
	if m.CreateBy == 0 {
		// 从context中获取createBy
		if createBy, ok := tx.Statement.Context.Value(UserIdContextKey).(int64); ok {
			m.CreateBy = createBy
		}
	}
	return nil
}

type CreateMCPServerReq struct {
	Name        string `json:"name" binding:"required"`
	Url         string `json:"url" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateMCPServerReq struct {
	Id          int64  `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Url         string `json:"url" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type MCPToolSet struct {
	Name  string      `json:"name"`
	Tools []*mcp.Tool `json:"tools"`
}
