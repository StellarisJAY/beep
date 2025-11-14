package types

import (
	"github.com/cloudwego/eino/components/model"
)

// Chunk 文档切片
// 文档切片是文档的一个片段，用于存储在向量数据库中
type Chunk interface {
	Content() string             // 切片内容
	DocId() int64                // 文档ID
	Id() string                  // 切片ID
	DenseVector() []float32      // 稠密向量
	Metadata() map[string]string // 元数据
}

type ListChunksQuery struct {
	DocId    string `json:"docId" form:"docId"`
	Page     bool   `json:"page" form:"page"`
	PageNum  int64  `json:"pageNum" form:"pageNum"`
	PageSize int64  `json:"pageSize" form:"pageSize"`
}

type QueriedChunk struct {
	Id            string  `json:"id"`            // 主键ID
	DocId         string  `json:"docId"`         // 文档ID
	SliceId       string  `json:"sliceId"`       // 切片ID
	Content       string  `json:"content"`       // 切片内容
	Score         float64 `json:"score"`         // 相似度得分
	VectorScore   float64 `json:"vectorScore"`   // 向量相似度得分
	FulltextScore float64 `json:"fulltextScore"` // 全文搜索相似度得分
}

type SearchType string

const (
	SearchTypeFulltext SearchType = "fulltext"
	SearchTypeVector   SearchType = "vector"
	SearchTypeHybrid   SearchType = "hybrid"
)

type HybridType string

const (
	HybridTypeWeight HybridType = "weight"
	HybridTypeRerank HybridType = "rerank"
)

type SearchReq struct {
	Text       string     // 查询文本
	TopK       int        // TopK
	Type       SearchType // 搜索类型: fulltext, vector, hybrid
	Threshold  float64    // 相似度阈值
	HybridType HybridType // 混合搜索类型: weight, rerank
	Weight     float64    // 混合搜索向量权重
	Embedding  []float32  // 查询文本向量
}

// ParseInfo 文档解析信息
type ParseInfo struct {
	Content              []byte              // 文本内容
	DocId                string              // 文档ID
	KbId                 string              // 知识库ID
	ChunkOptions         ChunkOptions        // 切片选项
	ChatModel            model.BaseChatModel // 聊天模型
	EnableKnowledgeGraph bool                // 是否开启知识图谱
	OriginalFileName     string              // 原始文件名
	FileType             string              // 文件类型
	UseOcr               bool                // 使用OCR识别图像
}
