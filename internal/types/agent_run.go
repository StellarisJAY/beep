package types

const (
	AgentRunModeNormal = "normal"
	AgentRunModeDebug  = "debug"
)

type AgentRunReq struct {
	AgentId        int64  `json:"agent_id"`
	Agent          *Agent `json:"agent"`
	ConversationId int64  `json:"conversation_id"`
	Query          string `json:"query"`
	Mode           string `json:"mode"`
}

type AgentRunResp struct {
	MessageChan chan Message
	ErrorChan   chan error
}
