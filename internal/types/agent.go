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

// Agent 智能体
type Agent struct {
	BaseEntity
	Name         string      `json:"name" gorm:"not null;type:varchar(64);"`           // 智能体名称
	Description  string      `json:"description" gorm:"not null;type:varchar(255);"`   // 智能体描述
	Type         AgentType   `json:"type" gorm:"not null;"`                            // 智能体类型
	Config       AgentConfig `json:"config" gorm:"not null;type:json;"`                // 智能体配置
	WorkspaceId  string      `json:"workspace_id" gorm:"not null;type:varchar(36);"`   // 工作空间ID
	CreateBy     string      `json:"create_by" gorm:"not null;type:varchar(36);"`      // 创建人ID
	LastUpdateBy string      `json:"last_update_by" gorm:"not null;type:varchar(36);"` // 最后更新人ID
}

func (*Agent) TableName() string {
	return "agents"
}

func (a *Agent) BeforeCreate(tx *gorm.DB) error {
	if err := a.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if a.WorkspaceId == "" {
		// 从context中获取workspaceId
		if workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(string); ok {
			a.WorkspaceId = workspaceId
		}
	}
	if a.CreateBy == "" {
		if userId, ok := tx.Statement.Context.Value(UserIdContextKey).(string); ok {
			a.CreateBy = userId
		}
	}
	return nil
}

func (a *Agent) ChatModelId() string {
	// 从配置中获取聊天模型ID
	return a.Config.ReAct.ChatModel
}

func (a *Agent) McpServerIds() []string {
	if a.Config.ReAct == nil {
		return nil
	}
	return a.Config.ReAct.McpServers
}

func (a *Agent) KnowledgeBaseIds() []string {
	if a.Config.ReAct == nil {
		return nil
	}
	return a.Config.ReAct.KnowledgeBases
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
	Id          string       `json:"id" binding:"required"`
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
	ChatModel       string          `json:"chat_model"`       // 聊天模型ID
	KnowledgeBases  []string        `json:"knowledge_bases"`  // 关联的知识库ID列表
	McpServers      []string        `json:"mcp_servers"`      // 关联的MCP服务器ID列表
	Prompt          string          `json:"prompt"`           // 智能体的提示词
	MaxIterations   int             `json:"max_iterations"`   // 最大迭代次数
	MemoryOption    MemoryOption    `json:"memory_option"`    // 记忆选项
	RetrieverOption RetrieverOption `json:"retriever_option"` // 知识库检索选项
	UseSystemTools  bool            `json:"use_system_tools"` // 是否使用系统工具集
}

// WorkflowAgentConfig 是 Workflow 类型的智能体配置
type WorkflowAgentConfig struct {
}

type AgentDetail struct {
	Agent
	McpServers     []*MCPServer     `json:"mcp_servers"`     // 关联的MCP服务器列表
	KnowledgeBases []*KnowledgeBase `json:"knowledge_bases"` // 关联的知识库列表
}
