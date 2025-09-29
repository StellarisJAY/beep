package interfaces

import (
	"beep/internal/types"
	"context"
)

type AgentRun interface {
	Run(ctx context.Context, req types.AgentRunReq) (*types.AgentRunResp, error)
	Cancel(ctx context.Context)
}
