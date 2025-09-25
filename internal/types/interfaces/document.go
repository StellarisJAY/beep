package interfaces

import (
	"beep/internal/types"
	"context"
	"mime/multipart"
)

// DocumentRepo 文档数据库
type DocumentRepo interface {
	Create(ctx context.Context, document *types.Document) error
	Update(ctx context.Context, document *types.Document) error
	Get(ctx context.Context, id int64) (*types.Document, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, query types.DocumentQuery) ([]*types.Document, int, error)
	DeleteByKnowledgeBaseId(ctx context.Context, knowledgeBaseId int64) error
}

// DocumentService 文档服务
type DocumentService interface {
	// CreateFromFile 从上传文件创建文档
	CreateFromFile(ctx context.Context, knowledgeBaseId int64, file *multipart.FileHeader) error
	// Delete 删除文档
	Delete(ctx context.Context, id int64) error
	// List 文档列表
	List(ctx context.Context, query types.DocumentQuery) ([]*types.Document, int, error)
	// Rename 重命名文档
	Rename(ctx context.Context, req types.RenameDocumentReq) error
	// Download 下载文档
	Download(ctx context.Context, id int64) (string, error)
	// Parse 解析文档
	Parse(ctx context.Context, id int64) error
}
