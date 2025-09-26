package reader

import (
	"context"
)

type ImageReader struct {
}

// Read OCR读取图片文件内容
func (i ImageReader) Read(ctx context.Context, content []byte, info Info) (textContent string, err error) {
	return
}
