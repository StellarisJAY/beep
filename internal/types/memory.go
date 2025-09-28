package types

// ShortTermMemoryQuery 短期记忆查询
type ShortTermMemoryQuery struct {
	ConversationId int64  `json:"conversation_id"`
	WindowSize     int    `json:"window_size"`
	Keyword        string `json:"keyword"`
}
