package common

import (
	"beep/internal/models/chat"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
)

func GetStaticMemory(ctx context.Context, service interfaces.MemoryService, conversationId string, windowSize int) ([]*chat.Message, error) {
	query := types.ShortTermMemoryQuery{
		ConversationId: conversationId,
		WindowSize:     windowSize,
	}
	messages, err := service.SearchShortTermMemory(ctx, query)
	if err != nil {
		return nil, err
	}

	// 转换为 schema.Message 格式
	schemaMessages := make([]*chat.Message, 0, len(messages))
	for _, msg := range messages {
		schemaMessages = append(schemaMessages, convertMessage(msg))
	}
	return schemaMessages, nil
}
