package interfaces

import (
	"beep/internal/types"
	"context"
)

// ModelFactoryRepo 模型供应商数据库
type ModelFactoryRepo interface {
	// Create 创建模型供应商
	Create(ctx context.Context, mf *types.ModelFactory) error
	// Update 更新模型供应商
	Update(ctx context.Context, mf *types.ModelFactory) error
	// List 所有模型供应商
	List(ctx context.Context) ([]*types.ModelFactory, error)
	// Delete 删除模型供应商
	Delete(ctx context.Context, id string) error
}

// ModelRepo 模型数据库
type ModelRepo interface {
	// Create 创建模型
	Create(ctx context.Context, mf *types.Model) error
	// CreateMany 创建多个模型
	CreateMany(ctx context.Context, mfs []*types.Model) error
	// Update 更新模型
	Update(ctx context.Context, mf *types.Model) error
	// List 所有模型
	List(ctx context.Context, query types.ListModelQuery) ([]*types.Model, error)
	// Delete 删除模型
	Delete(ctx context.Context, id string) error
	// GetDetail 获取模型详情
	GetDetail(ctx context.Context, id string) (*types.ModelDetail, error)
}

// ModelService 模型服务
type ModelService interface {
	// CreateFactory 新增模型供应商
	CreateFactory(ctx context.Context, req types.CreateModelFactoryReq) error
	// UpdateFactory 更新模型供应商
	UpdateFactory(ctx context.Context, req types.UpdateModelFactoryReq) error
	// ListFactory 列表模型供应商
	ListFactory(ctx context.Context) ([]*types.ModelFactory, error)
	// ListModels 列表模型
	ListModels(ctx context.Context, query types.ListModelQuery) ([]*types.Model, error)
	// GetModelDetail 获取模型详情
	GetModelDetail(ctx context.Context, id string) (*types.ModelDetail, error)
}
