package interfaces

import (
	"beep/internal/types"
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPServerRepo MCP 服务器的数据库接口
type MCPServerRepo interface {
	// Create 创建 MCP 服务器
	Create(ctx context.Context, ms *types.MCPServer) error
	// Update 更新 MCP 服务器
	Update(ctx context.Context, ms *types.MCPServer) error
	// Delete 删除 MCP 服务器
	Delete(ctx context.Context, id int64) error
	// List 所有 MCP 服务器
	List(ctx context.Context) ([]*types.MCPServer, error)
	// Get 根据 ID 获取 MCP 服务器
	Get(ctx context.Context, id int64) (*types.MCPServer, error)
}

type MCPServerService interface {
	// Create 创建 MCP 服务器
	Create(ctx context.Context, req types.CreateMCPServerReq) error
	// Update 更新 MCP 服务器
	Update(ctx context.Context, req types.UpdateMCPServerReq) error
	// Delete 删除 MCP 服务器
	Delete(ctx context.Context, id int64) error
	// List 所有 MCP 服务器
	List(ctx context.Context) ([]*types.MCPServer, error)
	// Get 根据 ID 获取 MCP 服务器
	Get(ctx context.Context, id int64) (*types.MCPServer, error)
	// ListTools 获取 MCP 服务器工具列表，结果组装到 MCPServer 结构体的 Tools 字段
	ListTools(ctx context.Context, ms *types.MCPServer) error
	// Call 调用 MCP 服务器工具
	Call(ctx context.Context, id int64, request *mcp.CallToolParams) (*mcp.CallToolResult, error)
	CallWithElicitation(ctx context.Context, id int64, request *mcp.CallToolParams) (*mcp.CallToolResult, error)
}
