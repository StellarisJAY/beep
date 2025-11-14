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
	EmbeddingModel string       `json:"embedding_model" gorm:"not null;"` // 知识库的嵌入模型
	ChatModel      string       `json:"chat_model" gorm:"not null;"`      // 知识库的聊天模型
	WorkspaceId    string       `json:"workspace_id" gorm:"not null;type:varchar(36);"`
	Public         bool         `json:"public" gorm:"not null;"`
	CreateBy       string       `json:"create_by" gorm:"not null;type:varchar(36);"`
	LastUpdateBy   string       `json:"last_update_by" gorm:"not null;type:varchar(36);"`
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
	if kb.WorkspaceId == "" {
		// 从context中获取workspaceId
		if workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(string); ok {
			kb.WorkspaceId = workspaceId
		}
	}
	if kb.CreateBy == "" {
		// 从context中获取createBy
		if createBy, ok := tx.Statement.Context.Value(UserIdContextKey).(string); ok {
			kb.CreateBy = createBy
		}
	}
	return nil
}

func (kb *KnowledgeBase) StorageBucketName() string {
	return fmt.Sprintf("beep-kb-%s", kb.ID)
}

func (kb *KnowledgeBase) StorageCollectionName() string {
	return fmt.Sprintf("beep-kb-%s", kb.ID)
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
	Ids        []string
}

type CreateKnowledgeBaseReq struct {
	Name           string       `json:"name" binding:"required"`
	Description    string       `json:"description" binding:"required"`
	EmbeddingModel string       `json:"embedding_model,string" binding:"required"`
	ChatModel      string       `json:"chat_model,string" binding:"required"`
	Public         *bool        `json:"public" binding:"required"`
	ChunkOptions   ChunkOptions `json:"chunk_options" binding:"required"`
}

type UpdateKnowledgeBaseReq struct {
	Id           string       `json:"id,string" binding:"required"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ChatModel    string       `json:"chat_model,string"`
	Public       *bool        `json:"public"`
	ChunkOptions ChunkOptions `json:"chunk_options"`
}

// RetrieverOption 知识库检索选项
type RetrieverOption struct {
	TopK       int        `json:"top_k"`       // TopK
	Threshold  float64    `json:"threshold"`   // 相似度阈值
	SearchType SearchType `json:"search_type"` // 搜索类型: fulltext, vector, hybrid
	HybridType HybridType `json:"hybrid_type"` // 混合搜索类型: weight, rerank
	Weight     float64    `json:"weight"`      // 混合搜索向量权重
	Reranker   string     `json:"reranker"`    // 重排模型ID
}
