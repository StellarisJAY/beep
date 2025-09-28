package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type ConversationRepo struct {
	db *gorm.DB
}

func NewConversationRepo(db *gorm.DB) interfaces.ConversationRepo {
	return &ConversationRepo{db: db}
}

func (c *ConversationRepo) Create(ctx context.Context, conversation *types.Conversation) error {
	return c.db.WithContext(ctx).Create(conversation).Error
}

func (c *ConversationRepo) List(ctx context.Context, query types.ConversationQuery) ([]*types.Conversation, int, error) {
	var conversations []*types.Conversation
	var total int64
	d := c.db.WithContext(ctx).Model(&types.Conversation{})
	if query.UserId != 0 {
		d = d.Where("user_id = ?", query.UserId)
	}
	if query.Title != "" {
		d = d.Where("title like ?", "%"+query.Title+"%")
	}
	if err := d.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.Scopes(pageScope(query.Paged, query.PageNum, query.PageSize)).Find(&conversations).Error; err != nil {
		return nil, 0, err
	}
	return conversations, int(total), nil
}

func (c *ConversationRepo) FindById(ctx context.Context, id int64) (*types.Conversation, error) {
	conversation := &types.Conversation{}
	if err := c.db.WithContext(ctx).First(conversation, id).Error; err != nil {
		return nil, err
	}
	return conversation, nil
}

func (c *ConversationRepo) Delete(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Delete(&types.Conversation{}, "id = ?", id).Error
}

type MessageRepo struct {
	db *gorm.DB
}

func (m *MessageRepo) Create(ctx context.Context, message *types.Message) error {
	return m.db.WithContext(ctx).Create(message).Error
}

func (m *MessageRepo) List(ctx context.Context, conversationId int64) ([]*types.Message, error) {
	var messages []*types.Message
	if err := m.db.WithContext(ctx).Where("conversation_id = ?", conversationId).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *MessageRepo) Search(ctx context.Context, query types.MessageQuery) ([]*types.Message, error) {
	var messages []*types.Message
	d := m.db.WithContext(ctx).Model(&types.Message{}).Where("conversation_id = ?", query.ConversationId)
	// TODO 全文搜索
	if query.Keyword != "" {
		d = d.Where("content LIKE ?", "%"+query.Keyword+"%")
	}
	if err := d.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *MessageRepo) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&types.Message{}, "id = ?", id).Error
}

func NewMessageRepo(db *gorm.DB) interfaces.MessageRepo {
	return &MessageRepo{db: db}
}
