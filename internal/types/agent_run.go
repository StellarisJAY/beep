package types

type AgentRunReq struct {
	Agent          *Agent
	ConversationId int64
	Query          string
}

type AgentRunResp struct {
	MessageChan chan Message
	ErrorChan   chan error
}
