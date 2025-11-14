package types

const (
	ContextKeySSEWriter = "sse_writer"
)

// SendMessageReq 发送消息请求
type SendMessageReq struct {
	AgentID        string `json:"agent_id"` // 智能体ID
	Agent          *Agent `json:"agent"`    // 智能体信息
	ChatModelID    string `json:"chat_model_id"`
	ConversationID string `json:"conversation_id"` // 会话ID
	Query          string `json:"query"`           // 查询内容
}
