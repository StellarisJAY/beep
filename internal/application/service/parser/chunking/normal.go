package chunking

import (
	"beep/internal/types"
	"context"
)

type NormalStrategy struct {
}

// Chunk 普通分块策略
func (n NormalStrategy) Chunk(ctx context.Context, textContent string, options types.ChunkOptions) (chunks []types.Chunk, err error) {
	return
}
