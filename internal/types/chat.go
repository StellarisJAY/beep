package types

const (
	ContextKeySSEWriter = "sse_writer"
)

// SendMessageReq 发送消息请求
type SendMessageReq struct {
	AgentID        int64  `json:"agent_id,string"` // 智能体ID
	Agent          *Agent `json:"agent"`           // 智能体信息
	ChatModelID    int64  `json:"chat_model_id,string"`
	ConversationID int64  `json:"conversation_id,string"` // 会话ID
	Query          string `json:"query"`                  // 查询内容
}
