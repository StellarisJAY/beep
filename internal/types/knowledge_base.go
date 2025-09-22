package types

// KnowledgeBase 知识库
type KnowledgeBase struct {
	BaseEntity
	Name           string `json:"name" gorm:"not null;type:varchar(64);"`
	Description    string `json:"description"  gorm:"not null;type:varchar(255);"`
	EmbeddingModel int64  `json:"embedding_model" gorm:"not null;"`
	ChatModel      int64  `json:"chat_model" gorm:"not null;"`
	WorkspaceId    int64  `json:"workspaceId" gorm:"not null;"`
	CreateBy       int64  `json:"createBy"`
	LastUpdateBy   int64  `json:"lastUpdateBy"`
}

func (*KnowledgeBase) TableName() string {
	return "knowledge_bases"
}
