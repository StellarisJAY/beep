package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type AgentType string

const (
	AgentTypeReAct    AgentType = "react"
	AgentTypeWorkflow AgentType = "workflow"
)

type Agent struct {
	BaseEntity
	Name         string      `json:"name" gorm:"not null;type:varchar(64);"`
	Description  string      `json:"description" gorm:"not null;type:varchar(255);"`
	Type         AgentType   `json:"type" gorm:"not null;"`
	Config       AgentConfig `json:"config" gorm:"not null;type:json;"`
	WorkspaceId  int64       `json:"workspace_id" gorm:"not null;"`
	CreateBy     int64       `json:"create_by"`
	LastUpdateBy int64       `json:"last_update_by"`
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
	if a.CreateBy == 0 {
		if userId, ok := tx.Statement.Context.Value(UserIdContextKey).(int64); ok {
			a.CreateBy = userId
		}
	}
	return nil
}

type AgentQuery struct {
	Name       string `form:"name"`
	Type       string `form:"type"`
	CreateByMe bool   `form:"create_by_me"`
}

type CreateAgentReq struct {
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description" binding:"required"`
	Type        AgentType    `json:"type" binding:"required"`
	Config      *AgentConfig `json:"config" binding:"required"`
}

type UpdateAgentReq struct {
	Id          int64        `json:"id" binding:"required"`
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description" binding:"required"`
	Config      *AgentConfig `json:"config" binding:"required"`
}

type AgentConfig struct {
	ReAct    *ReActAgentConfig    `json:"re_act"`
	Workflow *WorkflowAgentConfig `json:"workflow"`
}

func (a AgentConfig) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AgentConfig) Scan(value any) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, a)
}

type MemoryControlType string

const (
	MemoryControlAgent  MemoryControlType = "agent"  // 智能体自己控制记忆，通过MCP自主查询和添加
	MemoryControlStatic MemoryControlType = "static" // 每次对话前后自动查询和添加记忆
)

type MemoryOption struct {
	EnableShortTermMemory bool              `json:"enable_short_term_memory"` // 是否启用短期记忆
	EnableLongTermMemory  bool              `json:"enable_long_term_memory"`  // 是否启用长期记忆
	MemoryControl         MemoryControlType `json:"memory_control"`           // 记忆控制策略
	MemoryWindowSize      int64             `json:"memory_window_size"`       // 记忆窗口大小
}

// ReActAgentConfig 是 ReAct 类型的智能体配置
type ReActAgentConfig struct {
	ChatModel       int64           `json:"chat_model"`       // 聊天模型ID
	KnowledgeBases  []int64         `json:"knowledge_bases"`  // 关联的知识库ID列表
	McpServers      []int64         `json:"mcp_servers"`      // 关联的MCP服务器ID列表
	Prompt          string          `json:"prompt"`           // 智能体的提示词
	MaxIterations   int             `json:"max_iterations"`   // 最大迭代次数
	MemoryOption    MemoryOption    `json:"memory_option"`    // 记忆选项
	RetrieverOption RetrieverOption `json:"retriever_option"` // 知识库检索选项
}

// WorkflowAgentConfig 是 Workflow 类型的智能体配置
type WorkflowAgentConfig struct {
}
