package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// KnowledgeBase 知识库
type KnowledgeBase struct {
	BaseEntity
	Name           string       `json:"name" gorm:"not null;type:varchar(64);"`
	Description    string       `json:"description"  gorm:"not null;type:varchar(255);"`
	EmbeddingModel int64        `json:"embedding_model" gorm:"not null;"` // 知识库的嵌入模型
	ChatModel      int64        `json:"chat_model" gorm:"not null;"`      // 知识库的聊天模型
	WorkspaceId    int64        `json:"workspace_id" gorm:"not null;"`
	Public         bool         `json:"public" gorm:"not null;"`
	CreateBy       int64        `json:"create_by"`
	LastUpdateBy   int64        `json:"last_update_by"`
	ChunkOptions   ChunkOptions `json:"chunk_options" gorm:"type:json;"` // 切片选项
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

func (kb *KnowledgeBase) StorageBucketName() string {
	return fmt.Sprintf("beep-kb-%d", kb.ID)
}

func (kb *KnowledgeBase) StorageCollectionName() string {
	return fmt.Sprintf("beep-kb-%d", kb.ID)
}

// ChunkOptions 切片选项
type ChunkOptions struct {
	ChunkSize    int      `json:"chunk_size"`    // 切片大小
	ChunkOverlap int      `json:"chunk_overlap"` // 切片重叠大小
	Separators   []string `json:"separators"`    // 切片分隔符
}

func (c *ChunkOptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal ChunkOptions value")
	}
	return json.Unmarshal(b, &c)
}

func (c ChunkOptions) Value() (driver.Value, error) {
	return json.Marshal(&c)
}

type KnowledgeBaseQuery struct {
	BaseQuery
	Name       string `form:"name"`
	CreateByMe bool   `form:"create_by_me"`
}

type CreateKnowledgeBaseReq struct {
	Name           string       `json:"name" binding:"required"`
	Description    string       `json:"description" binding:"required"`
	EmbeddingModel int64        `json:"embedding_model,string" binding:"required"`
	ChatModel      int64        `json:"chat_model,string" binding:"required"`
	Public         *bool        `json:"public" binding:"required"`
	ChunkOptions   ChunkOptions `json:"chunk_options" binding:"required"`
}

type UpdateKnowledgeBaseReq struct {
	Id           int64        `json:"id,string" binding:"required"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ChatModel    int64        `json:"chat_model,string"`
	Public       *bool        `json:"public"`
	ChunkOptions ChunkOptions `json:"chunk_options"`
}
