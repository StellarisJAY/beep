package embedding

import (
	"beep/internal/types"
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Embedder 嵌入模型接口
type Embedder interface {
	// EmbedString 向量化一个字符串
	EmbedString(ctx context.Context, text string) ([]float64, error)
	// BatchEmbedStrings 批量向量化字符串
	BatchEmbedStrings(ctx context.Context, texts []string) ([][]float64, error)
}

// NewEmbedder 创建一个嵌入模型
func NewEmbedder(model types.ModelDetail) Embedder {
	switch model.FactoryType {
	case types.FactoryDashscope, types.FactoryOpenAI: // OpenAI-compatible 模型
		client := openai.NewClient(option.WithBaseURL(model.BaseUrl), option.WithAPIKey(model.ApiKeyDecrypted))
		return &remoteEmbedder{
			client:     client,
			modelName:  model.Name,
			dimensions: 1024, // TODO dimensions
		}
	case types.FactoryOllama:
		// TODO 实现ollama的embedder
		panic("not implemented")
	default:
		panic("not implemented")
	}
}
