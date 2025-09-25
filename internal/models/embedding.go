package models

import (
	"beep/internal/types"
	"context"
	"errors"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
)

// CreateEmbedder 创建嵌入模型
func CreateEmbedder(modelDetail types.ModelDetail) (embedding.Embedder, error) {
	switch modelDetail.FactoryType {
	case types.FactoryOpenAI:
		return openai.NewEmbedder(context.Background(), &openai.EmbeddingConfig{
			APIKey:  modelDetail.ApiKeyDecrypted,
			BaseURL: modelDetail.BaseUrl,
			Model:   modelDetail.Name,
		})
	case types.FactoryDashscope:
		return dashscope.NewEmbedder(context.Background(), &dashscope.EmbeddingConfig{
			APIKey: modelDetail.ApiKeyDecrypted,
			Model:  modelDetail.Name,
		})
	case types.FactoryOllama:
		return ollama.NewEmbedder(context.Background(), &ollama.EmbeddingConfig{
			BaseURL: modelDetail.BaseUrl,
			Model:   modelDetail.Name,
		})
	default:
		return nil, errors.New("unsupported factory type")
	}
}
