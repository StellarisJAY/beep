package models

import (
	"beep/internal/types"
	"context"
	"errors"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
)

// NewChatModel 创建聊天模型
func NewChatModel(modelDetail types.ModelDetail, option types.ChatModelOption) (model.BaseChatModel, error) {
	switch modelDetail.FactoryType {
	case types.FactoryOpenAI:
		maxTokens := int(modelDetail.MaxTokens)
		return openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
			APIKey:              modelDetail.ApiKeyDecrypted,
			Model:               modelDetail.Name,
			MaxCompletionTokens: &maxTokens,
			Temperature:         option.Temperature,
			TopP:                option.TopP,
			Stop:                option.Stop,
			PresencePenalty:     option.PresencePenalty,
			FrequencyPenalty:    option.FrequencyPenalty,
			Seed:                option.Seed,
		})
	case types.FactoryDashscope:
		maxTokens := int(modelDetail.MaxTokens)
		return qwen.NewChatModel(context.Background(), &qwen.ChatModelConfig{
			BaseURL:          modelDetail.BaseUrl,
			APIKey:           modelDetail.ApiKeyDecrypted,
			Model:            modelDetail.Name,
			MaxTokens:        &maxTokens,
			Temperature:      option.Temperature,
			TopP:             option.TopP,
			Stop:             option.Stop,
			PresencePenalty:  option.PresencePenalty,
			Seed:             option.Seed,
			FrequencyPenalty: option.FrequencyPenalty,
			EnableThinking:   option.Thinking,
		})
	default:
		return nil, errors.New("model factory not found")
	}
}
