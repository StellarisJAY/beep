package interfaces

import (
	"beep/internal/types"
	"context"
)

// KnowledgeBaseRepo 知识库数据库接口
type KnowledgeBaseRepo interface {
	// Create 创建知识库
	Create(ctx context.Context, kb *types.KnowledgeBase) error
	// List 知识库列表
	List(ctx context.Context, query types.KnowledgeBaseQuery) ([]*types.KnowledgeBase, int, error)
	// FindById 根据ID查找知识库
	FindById(ctx context.Context, id string) (*types.KnowledgeBase, error)
	// Update 更新知识库
	Update(ctx context.Context, kb *types.KnowledgeBase) error
	// Delete 删除知识库
	Delete(ctx context.Context, id string) error
}

// KnowledgeBaseService 知识库服务接口
type KnowledgeBaseService interface {
	// Create 创建知识库
	Create(ctx context.Context, req types.CreateKnowledgeBaseReq) error
	// List 知识库列表
	List(ctx context.Context, query types.KnowledgeBaseQuery) ([]*types.KnowledgeBase, int, error)
	// Update 更新知识库
	Update(ctx context.Context, req types.UpdateKnowledgeBaseReq) error
	// Delete 删除知识库
	Delete(ctx context.Context, id string) error
}
