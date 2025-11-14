package service

import (
	"beep/internal/application/repository"
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	errors2 "errors"
	"log/slog"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

type MCPServerService struct {
	mcpServerRepo interfaces.MCPServerRepo
}

func NewMCPServerService(db *gorm.DB) interfaces.MCPServerService {
	return &MCPServerService{
		mcpServerRepo: repository.NewMCPServerRepo(db),
	}
}

func (m *MCPServerService) Create(ctx context.Context, req types.CreateMCPServerReq) error {
	data := &types.MCPServer{Name: req.Name, Url: req.Url, Description: req.Description}
	if err := m.mcpServerRepo.Create(ctx, data); err != nil {
		return errors.NewInternalServerError("新增MCP服务失败", err)
	}
	return nil
}

func (m *MCPServerService) Update(ctx context.Context, req types.UpdateMCPServerReq) error {
	ms, _ := m.mcpServerRepo.Get(ctx, req.Id)
	if ms == nil {
		return errors.NewNotFoundError("MCP服务不存在", nil)
	}
	ms.Name = req.Name
	ms.Url = req.Url
	ms.Description = req.Description
	if err := m.mcpServerRepo.Update(ctx, ms); err != nil {
		return errors.NewInternalServerError("更新MCP服务失败", err)
	}
	return nil
}

func (m *MCPServerService) Delete(ctx context.Context, id int64) error {
	if err := m.mcpServerRepo.Delete(ctx, id); err != nil {
		return errors.NewInternalServerError("删除MCP服务失败", err)
	}
	// TODO 删除关联
	return nil
}

func (m *MCPServerService) List(ctx context.Context) ([]*types.MCPServer, error) {
	list, err := m.mcpServerRepo.List(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError("获取MCP服务列表失败", err)
	}
	return list, nil
}

func (m *MCPServerService) Get(ctx context.Context, id int64) (*types.MCPServer, error) {
	ms, err := m.mcpServerRepo.Get(ctx, id)
	if errors2.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NewNotFoundError("MCP服务不存在", nil)
	}
	if err != nil {
		return nil, errors.NewInternalServerError("获取MCP服务失败", err)
	}
	return ms, nil
}

func (m *MCPServerService) ListTools(ctx context.Context, ms *types.MCPServer) error {
	if ms.Url == "" || !strings.HasPrefix(ms.Url, "http://") && !strings.HasPrefix(ms.Url, "https://") {
		return errors.NewBadRequestError("MCP服务URL格式错误", nil)
	}
	cli := mcp.NewClient(&mcp.Implementation{Name: "mcp-cli", Version: "v1.0.0"}, nil)
	session, err := cli.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: ms.Url}, nil)
	if err != nil {
		return errors.NewInternalServerError("MCP服务连接失败", err)
	}
	defer func() {
		_ = session.Close()
	}()
	res, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return errors.NewInternalServerError("MCP服务获取工具列表失败", err)
	}
	ms.Tools = res.Tools
	ms.Available = true
	return nil
}

func (m *MCPServerService) getMCPServer(ctx context.Context, id int64) (*types.MCPServer, error) {
	mcpServer, err := m.mcpServerRepo.Get(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerError("获取MCP服务失败", err)
	}
	if mcpServer.Url == "" || !strings.HasPrefix(mcpServer.Url, "http://") && !strings.HasPrefix(mcpServer.Url, "https://") {
		return nil, errors.NewBadRequestError("MCP服务URL格式错误", nil)
	}
	return mcpServer, nil
}

func (m *MCPServerService) Call(ctx context.Context, id int64, request *mcp.CallToolParams) (*mcp.CallToolResult, error) {
	mcpServer, err := m.getMCPServer(ctx, id)
	if err != nil {
		return nil, err
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-cli", Version: "v1.0.0"}, &mcp.ClientOptions{})
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: mcpServer.Url}, nil)
	if err != nil {
		return nil, errors.NewInternalServerError("MCP服务连接失败", err)
	}
	defer func() {
		_ = session.Close()
	}()
	res, err := session.CallTool(ctx, request)
	if err != nil {
		return nil, errors.NewInternalServerError("MCP服务调用工具失败", err)
	}
	return res, nil
}

// CallWithElicitation 带有MCP征询机制的调用工具
func (m *MCPServerService) CallWithElicitation(ctx context.Context, id int64, request *mcp.CallToolParams) (*mcp.CallToolResult, error) {
	mcpServer, err := m.getMCPServer(ctx, id)
	if err != nil {
		return nil, err
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-cli", Version: "v1.0.0"}, &mcp.ClientOptions{
		ElicitationHandler: func(c context.Context, req *mcp.ElicitRequest) (*mcp.ElicitResult, error) {
			slog.Info("wait permission for", "tool", request.Name, "param", request.Arguments)
			// TODO 等待用户批准请求
			return &mcp.ElicitResult{Action: "accept"}, nil
		},
	})
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: mcpServer.Url}, nil)
	if err != nil {
		return nil, errors.NewInternalServerError("MCP服务连接失败", err)
	}
	defer func() {
		_ = session.Close()
	}()
	res, err := session.CallTool(ctx, request)
	if err != nil {
		return nil, errors.NewInternalServerError("MCP服务调用工具失败", err)
	}
	return res, nil
}

func (m *MCPServerService) ListWithoutTools(ctx context.Context, query types.MCPServerQuery) ([]*types.MCPServer, error) {
	list, err := m.mcpServerRepo.ListWithoutTools(ctx, query)
	if err != nil {
		return nil, errors.NewInternalServerError("获取MCP服务列表失败", err)
	}
	return list, nil
}
