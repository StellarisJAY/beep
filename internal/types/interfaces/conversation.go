package interfaces

import (
	"beep/internal/types"
	"context"
)

type ConversationRepo interface {
	// Create 创建对话
	Create(ctx context.Context, conversation *types.Conversation) error
	// List 查询对话列表
	List(ctx context.Context, query types.ConversationQuery) ([]*types.Conversation, int, error)
	// FindById 根据ID查询对话
	FindById(ctx context.Context, id int64) (*types.Conversation, error)
	// Delete 删除对话
	Delete(ctx context.Context, id int64) error
	// UpdateTitle 更新对话标题
	UpdateTitle(ctx context.Context, id int64, title string) error
}

type MessageRepo interface {
	// Create 创建消息
	Create(ctx context.Context, message *types.Message) error
	// List 查询消息列表
	List(ctx context.Context, conversationId int64) ([]*types.Message, error)
	// Search 搜索消息
	Search(ctx context.Context, query types.MessageQuery) ([]*types.Message, error)
	// Delete 删除消息
	Delete(ctx context.Context, id int64) error
	// Update 更新消息
	Update(ctx context.Context, message *types.Message) error
}

type ConversationService interface {
	List(ctx context.Context, query types.ConversationQuery) ([]*types.Conversation, error)
	ListMessages(ctx context.Context, conversationId int64) ([]*types.Message, error)
}
