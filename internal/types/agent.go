package types

import "gorm.io/gorm"

type AgentType string

const (
	AgentTypeSimple   AgentType = "simple"
	AgentTypeWorkflow AgentType = "workflow"
)

type Agent struct {
	BaseEntity
	Name         string    `json:"name" gorm:"not null;type:varchar(64);"`
	Description  string    `json:"description" gorm:"not null;type:varchar(255);"`
	Type         AgentType `json:"type" gorm:"not null;"`
	Config       string    `json:"config" gorm:"not null;type:text;"`
	WorkspaceId  int64     `json:"workspace_id" gorm:"not null;"`
	CreateBy     int64     `json:"create_by"`
	LastUpdateBy int64     `json:"last_update_by"`
}

func (*Agent) TableName() string {
	return "agents"
}

func (a *Agent) BeforeCreate(tx *gorm.DB) error {
	if err := a.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if a.WorkspaceId == 0 {
		// 从context中获取workspaceId
		if workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(int64); ok {
			a.WorkspaceId = workspaceId
		}
	}
	return nil
}
