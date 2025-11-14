package types

import "gorm.io/gorm"

const (
	ToolStatusNone   = 0
	ToolStatusWait   = 1 // 等待中
	ToolStatusAccept = 2 // 接受
	ToolStatusReject = 3 // 拒绝
)

// Message 消息
type Message struct {
	BaseEntity
	ConversationId string `json:"conversation_id" gorm:"index;not null;type:varchar(36);"` // 对话ID
	Role           string `json:"role" gorm:"type:varchar(16);not null;"`                  // 消息角色，user、system、assistant
	Content        string `json:"content" gorm:"type: text; not null;"`                    // 消息内容
	ToolCall       string `json:"tool_call" gorm:"type: text;"`                            // 工具调用
	ToolCallParams string `json:"tool_call_params" gorm:"type: text;"`                     // 工具调用参数
	ToolCallId     string `json:"tool_call_id" gorm:"type: varchar(255);"`                 // 工具调用ID
	ToolStatus     int    `json:"tool_status" gorm:"type:smallint;default:0;"`             // 工具状态，0：默认，1：等待中，2：接受，3：拒绝
}

func (*Message) TableName() string {
	return "messages"
}

// Conversation 对话
type Conversation struct {
	BaseEntity
	ModelId     string `json:"model_id" gorm:"not null;type:varchar(36);"`     // 模型ID，如果为null，表示对话是智能体对话
	AgentId     string `json:"agent_id" gorm:"not null;type:varchar(36);"`     // 智能体ID，如果为null，表示对话是模型对话
	Title       string `json:"title" gorm:"type:varchar(256);not null;"`       // 对话标题
	Summary     string `json:"summary" gorm:"type:text;not null;"`             // 对话摘要
	CreateBy    string `json:"create_by" gorm:"not null;type:varchar(36);"`    // 创建人ID
	WorkspaceId string `json:"workspace_id" gorm:"not null;type:varchar(36);"` // 工作空间ID
}

func (*Conversation) TableName() string {
	return "conversations"
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if err := c.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if c.CreateBy == "" {
		createBy, ok := tx.Statement.Context.Value(UserIdContextKey).(string)
		if ok {
			c.CreateBy = createBy
		}
	}
	if c.WorkspaceId == "" {
		workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(string)
		if ok {
			c.WorkspaceId = workspaceId
		}
	}
	return nil
}

type ConversationQuery struct {
	BaseQuery
	UserId      string `json:"user_id" form:"user_id"`
	Title       string `json:"title" form:"title"`
	GroupByDate bool   `json:"group_by_date" form:"group_by_date"` // 是否按最近、最近一周分组
}

type MessageQuery struct {
	ConversationId string `json:"conversation_id" form:"conversation_id"` // 对话ID
	Keyword        string `json:"keyword" form:"keyword"`                 // 搜索关键词
	LastN          int    `json:"last_n" form:"last_n"`
	Role           string `json:"role" form:"role"`           // 消息角色，user、system、assistant
	ToolCall       *bool  `json:"tool_call" form:"tool_call"` // 是否查询工具调用消息
}
