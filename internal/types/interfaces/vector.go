package interfaces

import (
	"beep/internal/types"
	"context"
)

// VectorStore 向量数据库
type VectorStore interface {
	Ping() error                                                                                                   // 检查连接
	Close() error                                                                                                  // 关闭连接
	CreateCollection(ctx context.Context, name string, denseDims int64) error                                      // 创建集合
	DropCollection(ctx context.Context, name string) error                                                         // 删除集合
	Add(ctx context.Context, coll string, chunks []types.Chunk) error                                              // 添加切片
	Delete(ctx context.Context, coll string, chunks []types.Chunk) error                                           // 删除切片
	ListChunks(ctx context.Context, coll string, query types.ListChunksQuery) ([]types.QueriedChunk, int64, error) // 列表切片
	Search(ctx context.Context, coll string, req types.SearchReq) ([]types.QueriedChunk, error)
}

// ParseService 解析服务
type ParseService interface {
	Parse(ctx context.Context, info types.ParseInfo) error // 解析文档
	CancelParse(ctx context.Context, docId int64) error    // 取消解析文档
}
