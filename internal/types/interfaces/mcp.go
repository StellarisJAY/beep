package interfaces

import (
	"beep/internal/types"
	"context"
)

// MCPServerRepo MCP 服务器的数据库接口
type MCPServerRepo interface {
	Create(ctx context.Context, ms *types.MCPServer) error
	Update(ctx context.Context, ms *types.MCPServer) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*types.MCPServer, error)
	Get(ctx context.Context, id int64) (*types.MCPServer, error)
}

type MCPServerService interface {
	Create(ctx context.Context, req types.CreateMCPServerReq) error
	Update(ctx context.Context, req types.UpdateMCPServerReq) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*types.MCPServer, error)
	Get(ctx context.Context, id int64) (*types.MCPServer, error)
	ListTools(ctx context.Context, ms *types.MCPServer) error
}
