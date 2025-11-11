package interfaces

import (
	"beep/internal/types"
	"context"
)

// AgentRun 智能体运行接口
type AgentRun interface {
	// Run 运行智能体
	Run(ctx context.Context, req types.AgentRunReq) (*types.AgentRunResp, error)
	// Cancel 取消智能体运行
	Cancel(ctx context.Context)
	// SignalTool 工具调用暂停时，发送通过或拒绝信号
	SignalTool(ctx context.Context, signal types.ToolSignal) (*types.AgentRunResp, error)
}

// AgentRunFactory 智能体运行工厂接口，由工厂创建智能体运行实例
type AgentRunFactory interface {
	// CreateAgentRun 创建智能体运行实例
	CreateAgentRun(req types.AgentRunReq) (AgentRun, error)
}
