package chunking

import (
	"beep/internal/types"
	"context"
)

type Strategy interface {
	Chunk(ctx context.Context, content string, options types.ChunkOptions) ([]types.Chunk, error)
}
