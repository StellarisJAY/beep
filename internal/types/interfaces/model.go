package interfaces

import (
	"beep/internal/types"
	"context"
)

type ModelFactoryRepo interface {
	Create(ctx context.Context, mf *types.ModelFactory) error
	Update(ctx context.Context, mf *types.ModelFactory) error
	List(ctx context.Context) ([]*types.ModelFactory, error)
	Delete(ctx context.Context, id int64) error
}

type ModelRepo interface {
	Create(ctx context.Context, mf *types.Model) error
	CreateMany(ctx context.Context, mfs []*types.Model) error
	Update(ctx context.Context, mf *types.Model) error
	List(ctx context.Context, query types.ListModelQuery) ([]*types.Model, error)
	Delete(ctx context.Context, id int64) error
	GetDetail(ctx context.Context, id int64) (*types.ModelDetail, error)
}

type ModelService interface {
	// CreateFactory 新增模型工厂
	CreateFactory(ctx context.Context, req types.CreateModelFactoryReq) error
	// UpdateFactory 更新模型工厂
	UpdateFactory(ctx context.Context, req types.UpdateModelFactoryReq) error
	// ListFactory 列表模型工厂
	ListFactory(ctx context.Context) ([]*types.ModelFactory, error)
	// ListModels 列表模型
	ListModels(ctx context.Context, query types.ListModelQuery) ([]*types.Model, error)
	// GetModelDetail 获取模型详情
	GetModelDetail(ctx context.Context, id int64) (*types.ModelDetail, error)
}
