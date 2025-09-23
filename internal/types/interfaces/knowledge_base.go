package interfaces

import (
	"beep/internal/types"
	"context"
)

type KnowledgeBaseRepo interface {
	Create(ctx context.Context, kb *types.KnowledgeBase) error
	List(ctx context.Context, query types.KnowledgeBaseQuery) ([]*types.KnowledgeBase, int, error)
	FindById(ctx context.Context, id int64) (*types.KnowledgeBase, error)
	Update(ctx context.Context, kb *types.KnowledgeBase) error
	Delete(ctx context.Context, id int64) error
}

type KnowledgeBaseService interface {
	Create(ctx context.Context, req types.CreateKnowledgeBaseReq) error
	List(ctx context.Context, query types.KnowledgeBaseQuery) ([]*types.KnowledgeBase, int, error)
	Update(ctx context.Context, req types.UpdateKnowledgeBaseReq) error
	Delete(ctx context.Context, id int64) error
}
