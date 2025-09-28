package service

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
)

type MemoryService struct {
	messageRepo      interfaces.MessageRepo
	conversationRepo interfaces.ConversationRepo
	vectorStore      interfaces.VectorStore
}

func NewMemoryService(messageRepo interfaces.MessageRepo,
	conversationRepo interfaces.ConversationRepo,
	vectorStore interfaces.VectorStore) interfaces.MemoryService {
	return &MemoryService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		vectorStore:      vectorStore,
	}
}

func (m *MemoryService) SearchShortTermMemory(ctx context.Context, query types.ShortTermMemoryQuery) ([]*types.Message, error) {
	messages, err := m.messageRepo.Search(ctx, types.MessageQuery{
		ConversationId: query.ConversationId,
		Keyword:        query.Keyword,
		LastN:          query.WindowSize,
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}
