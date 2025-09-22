package types

// KnowledgeBase 知识库
type KnowledgeBase struct {
	BaseEntity
	Name           string `json:"name"`
	Description    string `json:"description"`
	EmbeddingModel uint64 `json:"embedding_model"`
	ChatModel      uint64 `json:"chat_model"`
	WorkspaceId    uint64 `json:"workspaceId"`
	CreateBy       uint64 `json:"createBy"`
	LastUpdateBy   uint64 `json:"lastUpdateBy"`
}

func (*KnowledgeBase) TableName() string {
	return "knowledge_bases"
}
