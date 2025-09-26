package reader

import (
	"context"
	"errors"
)

type PdfReader struct {
}

// Read 读取PDF文件内容
func (p PdfReader) Read(ctx context.Context, content []byte, info Info) (string, error) {
	// 如果不使用OCR
	if !info.UseOcr {
		return "", errors.New("pdf reader not implemented")
	}
	return "", nil
}
