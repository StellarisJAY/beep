package parser

import (
	"beep/internal/application/service/parser/reader"
	"beep/internal/dag"
	"beep/internal/types"
	"context"
)

func (d *DocumentParser) chunking(ctx context.Context, node dag.Node) (result dag.NodeFuncReturn) {
	parseInfo := ctx.Value(parseInfoContextKey).(*types.ParseInfo)
	// 提取文档内容
	fr := reader.NewFileReader(parseInfo.FileType)
	textContent, err := fr.Read(ctx, parseInfo.Content, reader.Info{
		UseOcr:           parseInfo.UseOcr,
		FileType:         parseInfo.FileType,
		OriginalFileName: parseInfo.OriginalFileName,
	})
	if err != nil {
		result.Error = err
		return
	}
	// TODO 切片内容
	result.Output = map[string]any{
		"text_content": textContent,
	}
	return
}
