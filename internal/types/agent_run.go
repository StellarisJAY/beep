package types

const (
	AgentRunModeNormal = "normal"
	AgentRunModeDebug  = "debug"
)

type AgentRunReq struct {
	AgentId        string `json:"agent_id"`
	Agent          *Agent `json:"agent"`
	ConversationId string `json:"conversation_id"`
	Query          string `json:"query"`
	Mode           string `json:"mode"`
}

type ToolSignal struct {
	AgentId        string `json:"agent_id"`
	ConversationId string `json:"conversation_id"`
	Accept         bool   `json:"accept"`
}

type AgentRunResp struct {
	MessageChan chan Message
	ErrorChan   chan error
}
