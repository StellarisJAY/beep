package reader

import (
	"context"
)

type WordReader struct {
}

// Read 读取Word文件内容
func (w WordReader) Read(ctx context.Context, content []byte, info Info) (textContent string, err error) {
	return
}
