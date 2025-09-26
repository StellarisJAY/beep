package reader

import (
	"context"
	"strings"
)

type Info struct {
	UseOcr           bool   // 使用OCR识别图像
	FileType         string // 文件类型
	OriginalFileName string // 原始文件名
}

type FileReader interface {
	// Read 读取文件内容
	Read(ctx context.Context, content []byte, info Info) (textContent string, err error)
}

func NewFileReader(fileType string) FileReader {
	switch {
	case strings.HasPrefix(fileType, "text/"):
		return TextReader{}
	case strings.HasPrefix(fileType, "image/"):
		return ImageReader{}
	case strings.HasPrefix(fileType, "application/pdf"):
		return PdfReader{}
	case strings.HasPrefix(fileType, "application/msword"):
		return WordReader{}
	default:
		return TextReader{}
	}
}
