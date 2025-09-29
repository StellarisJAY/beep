package interfaces

import (
	"beep/internal/types"
	"context"
)

type AgentRepo interface {
	Create(ctx context.Context, agent *types.Agent) error
	FindById(ctx context.Context, id int64) (*types.Agent, error)
	Update(ctx context.Context, agent *types.Agent) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, query types.AgentQuery) ([]*types.Agent, error)
}

type AgentService interface {
	Create(ctx context.Context, req types.CreateAgentReq) error
	Detail(ctx context.Context, id int64) (*types.Agent, error)
	Update(ctx context.Context, req types.UpdateAgentReq) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, query types.AgentQuery) ([]*types.Agent, error)
	Run(ctx context.Context, id int64, query string) error
}
