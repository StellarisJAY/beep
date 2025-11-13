package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
)

type ConversationService struct {
	conversationRepo interfaces.ConversationRepo
	messageRepo      interfaces.MessageRepo
}

func (c *ConversationService) List(ctx context.Context, query types.ConversationQuery) ([]*types.Conversation, error) {
	query.UserId = ctx.Value(types.UserIdContextKey).(int64)
	conversations, _, err := c.conversationRepo.List(ctx, query)
	if err != nil {
		return nil, errors.NewInternalServerError("查询对话列表失败", err)
	}
	return conversations, nil
}

func (c *ConversationService) ListMessages(ctx context.Context, conversationId int64) ([]*types.Message, error) {
	messages, err := c.messageRepo.List(ctx, conversationId)
	if err != nil {
		return nil, errors.NewInternalServerError("查询消息列表失败", err)
	}
	return messages, nil
}

func NewConversationService(conversationRepo interfaces.ConversationRepo, messageRepo interfaces.MessageRepo) interfaces.ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}
