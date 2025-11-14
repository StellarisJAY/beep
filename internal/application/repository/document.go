package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type DocumentRepo struct {
	db *gorm.DB
}

func NewDocumentRepo(db *gorm.DB) interfaces.DocumentRepo {
	return &DocumentRepo{db: db}
}

func (d *DocumentRepo) Create(ctx context.Context, document *types.Document) error {
	return d.db.WithContext(ctx).Model(&types.Document{}).Create(document).Error
}

func (d *DocumentRepo) Update(ctx context.Context, document *types.Document) error {
	return d.db.WithContext(ctx).Model(document).Where("id = ?", document.ID).Updates(document).Error
}

func (d *DocumentRepo) Get(ctx context.Context, id string) (*types.Document, error) {
	var document types.Document
	if err := d.db.WithContext(ctx).Scopes(workspaceScope(ctx)).First(&document, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (d *DocumentRepo) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&types.Document{}, "id = ?", id).Error
}

func (d *DocumentRepo) List(ctx context.Context, query types.DocumentQuery) ([]*types.Document, int, error) {
	var documents []*types.Document
	var total int64
	tx := d.db.WithContext(ctx).Model(&types.Document{}).
		Scopes(workspaceScope(ctx)).
		Where("knowledge_base_id = ?", query.KnowledgeBaseId)
	if query.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.ParseStatus != "" {
		tx = tx.Where("parse_status = ?", query.ParseStatus)
	}
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Scopes(pageScope(query.Paged, query.PageNum, query.PageSize)).Find(&documents).Error; err != nil {
		return nil, 0, err
	}
	return documents, int(total), nil
}

func (d *DocumentRepo) DeleteByKnowledgeBaseId(ctx context.Context, knowledgeBaseId string) error {
	return d.db.WithContext(ctx).Scopes(workspaceScope(ctx)).Where("knowledge_base_id = ?", knowledgeBaseId).Delete(&types.Document{}).Error
}
