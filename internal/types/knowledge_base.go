package types

import "gorm.io/gorm"

// KnowledgeBase 知识库
type KnowledgeBase struct {
	BaseEntity
	Name           string `json:"name" gorm:"not null;type:varchar(64);"`
	Description    string `json:"description"  gorm:"not null;type:varchar(255);"`
	EmbeddingModel int64  `json:"embedding_model" gorm:"not null;"`
	ChatModel      int64  `json:"chat_model" gorm:"not null;"`
	WorkspaceId    int64  `json:"workspaceId" gorm:"not null;"`
	Public         bool   `json:"public" gorm:"not null;"`
	CreateBy       int64  `json:"createBy"`
	LastUpdateBy   int64  `json:"lastUpdateBy"`
}

func (*KnowledgeBase) TableName() string {
	return "knowledge_bases"
}

// BeforeCreate 数据库插入前，自动设置id和workspace_id
func (kb *KnowledgeBase) BeforeCreate(tx *gorm.DB) error {
	if err := kb.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if kb.WorkspaceId == 0 {
		// 从context中获取workspaceId
		if workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(int64); ok {
			kb.WorkspaceId = workspaceId
		}
	}
	if kb.CreateBy == 0 {
		// 从context中获取createBy
		if createBy, ok := tx.Statement.Context.Value(UserIdContextKey).(int64); ok {
			kb.CreateBy = createBy
		}
	}
	return nil
}

type KnowledgeBaseQuery struct {
	BaseQuery
	Name       string `form:"name"`
	CreateByMe bool   `form:"create_by_me"`
}

type CreateKnowledgeBaseReq struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	EmbeddingModel int64  `json:"embedding_model,string" binding:"required"`
	ChatModel      int64  `json:"chat_model,string" binding:"required"`
	Public         *bool  `json:"public" binding:"required"`
}

type UpdateKnowledgeBaseReq struct {
	Id          int64  `json:"id,string" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ChatModel   int64  `json:"chat_model,string"`
	Public      *bool  `json:"public"`
}
