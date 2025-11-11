package types

import "gorm.io/gorm"

// Message 消息
type Message struct {
	BaseEntity
	ConversationId int64  `json:"conversation_id,string" gorm:"index;not null;"` // 对话ID
	Role           string `json:"role" gorm:"type:varchar(16);not null;"`        // 消息角色，user、system、assistant
	Content        string `json:"content" gorm:"type: text; not null;"`          // 消息内容
	ToolCall       string `json:"tool_call" gorm:"type: text;"`                  // 工具调用
	ToolCallParams string `json:"tool_call_params" gorm:"type: text;"`           // 工具调用参数
	ToolCallId     string `json:"tool_call_id" gorm:"type: varchar(255);"`       // 工具调用ID
}

func (*Message) TableName() string {
	return "messages"
}

// Conversation 对话
type Conversation struct {
	BaseEntity
	ModelId  int64  `json:"model_id,string"`                          // 模型ID，如果为null，表示对话是智能体对话
	AgentId  int64  `json:"agent_id,string"`                          // 智能体ID，如果为null，表示对话是模型对话
	Title    string `json:"title" gorm:"type:varchar(256);not null;"` // 对话标题
	Summary  string `json:"summary" gorm:"type:text;not null;"`       // 对话摘要
	CreateBy int64  `json:"create_by,string" gorm:"not null;"`        // 创建人ID
}

func (*Conversation) TableName() string {
	return "conversations"
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if err := c.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if c.CreateBy == 0 {
		createBy, ok := tx.Statement.Context.Value(UserIdContextKey).(int64)
		if ok {
			c.CreateBy = createBy
		}
	}
	return nil
}

type ConversationQuery struct {
	BaseQuery
	UserId      int64  `json:"user_id" form:"user_id"`
	Title       string `json:"title" form:"title"`
	GroupByDate bool   `json:"group_by_date" form:"group_by_date"` // 是否按最近、最近一周分组
}

type MessageQuery struct {
	ConversationId int64  `json:"conversation_id" form:"conversation_id"` // 对话ID
	Keyword        string `json:"keyword" form:"keyword"`                 // 搜索关键词
	LastN          int    `json:"last_n" form:"last_n"`
	Role           string `json:"role" form:"role"`           // 消息角色，user、system、assistant
	ToolCall       *bool  `json:"tool_call" form:"tool_call"` // 是否查询工具调用消息
}
