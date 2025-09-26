package reader

import (
	"context"
)

type TextReader struct{}

// Read 读取文本文件内容
func (t TextReader) Read(_ context.Context, content []byte, _ Info) (textContent string, err error) {
	textContent = string(content)
	return
}
