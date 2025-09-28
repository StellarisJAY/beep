package interfaces

import (
	"beep/internal/types"
	"context"
)

// MemoryService 智能体记忆服务
type MemoryService interface {
	// SearchShortTermMemory 搜索短期记忆
	SearchShortTermMemory(ctx context.Context, query types.ShortTermMemoryQuery) ([]*types.Message, error)
}
