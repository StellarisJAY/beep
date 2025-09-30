package interfaces

import (
	"beep/internal/types"
	"context"
)

// AgentRepo 智能体数据库接口
type AgentRepo interface {
	// Create 创建智能体
	Create(ctx context.Context, agent *types.Agent) error
	// FindById 根据id查询智能体
	FindById(ctx context.Context, id int64) (*types.Agent, error)
	// Update 更新智能体
	Update(ctx context.Context, agent *types.Agent) error
	// Delete 删除智能体
	Delete(ctx context.Context, id int64) error
	// List 查询智能体列表
	List(ctx context.Context, query types.AgentQuery) ([]*types.Agent, error)
}

// AgentService 智能体服务接口
type AgentService interface {
	// Create 创建智能体
	Create(ctx context.Context, req types.CreateAgentReq) error
	// Detail 查询智能体详情
	Detail(ctx context.Context, id int64) (*types.Agent, error)
	// Update 更新智能体
	Update(ctx context.Context, req types.UpdateAgentReq) error
	// Delete 删除智能体
	Delete(ctx context.Context, id int64) error
	// List 查询智能体列表
	List(ctx context.Context, query types.AgentQuery) ([]*types.Agent, error)
	// Run 运行智能体
	Run(ctx context.Context, req types.AgentRunReq) error
}
