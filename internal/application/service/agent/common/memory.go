package common

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"github.com/cloudwego/eino/schema"
)

func GetStaticMemory(ctx context.Context, service interfaces.MemoryService, conversationId int64, memoryOption types.MemoryOption) ([]*schema.Message, error) {
	query := types.ShortTermMemoryQuery{
		ConversationId: conversationId,
		WindowSize:     int(memoryOption.MemoryWindowSize),
	}
	messages, err := service.SearchShortTermMemory(context.Background(), query)
	if err != nil {
		return nil, err
	}

	// 转换为 schema.Message 格式
	schemaMessages := make([]*schema.Message, 0, len(messages))
	for _, msg := range messages {
		schemaMessages = append(schemaMessages, ConvertMessageToSchemaMessage(msg))
	}
	return schemaMessages, nil
}
