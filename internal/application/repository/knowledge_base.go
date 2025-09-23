package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type KnowledgeBaseRepo struct {
	db *gorm.DB
}

func NewKnowledgeBaseRepo(db *gorm.DB) interfaces.KnowledgeBaseRepo {
	return &KnowledgeBaseRepo{
		db: db,
	}
}

func (k *KnowledgeBaseRepo) Create(ctx context.Context, kb *types.KnowledgeBase) error {
	return k.db.WithContext(ctx).Model(kb).Create(kb).Error
}

func (k *KnowledgeBaseRepo) List(ctx context.Context, query types.KnowledgeBaseQuery) ([]*types.KnowledgeBase, int, error) {
	var list []*types.KnowledgeBase
	var total int64
	userId, _ := ctx.Value(types.UserIdContextKey).(int64)

	d := k.db.WithContext(ctx).Model(&types.KnowledgeBase{}).Scopes(workspaceScope(ctx))
	if query.Name != "" {
		d = d.Where("name = ?", query.Name)
	}
	if query.CreateByMe {
		d = d.Where("create_by = ?", userId)
	}
	if err := d.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := d.Scopes(pageScope(query.Paged, query.PageNum, query.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, int(total), nil
}

func (k *KnowledgeBaseRepo) FindById(ctx context.Context, id int64) (*types.KnowledgeBase, error) {
	var kb *types.KnowledgeBase
	err := k.db.WithContext(ctx).
		Model(&types.KnowledgeBase{}).
		Where("id = ?", id).
		Scopes(workspaceScope(ctx)).
		First(&kb).Error
	if err != nil {
		return nil, err
	}
	return kb, nil
}

func (k *KnowledgeBaseRepo) Update(ctx context.Context, kb *types.KnowledgeBase) error {
	return k.db.WithContext(ctx).Model(kb).Where("id = ?", kb.ID).Updates(kb).Error
}

func (k *KnowledgeBaseRepo) Delete(ctx context.Context, id int64) error {
	return k.db.WithContext(ctx).Delete(&types.KnowledgeBase{}, "id = ?", id).Error
}
