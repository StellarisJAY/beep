package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type ModelFactoryRepo struct {
	db *gorm.DB
}

func NewModelFactoryRepo(db *gorm.DB) interfaces.ModelFactoryRepo {
	return &ModelFactoryRepo{db: db}
}

func (m *ModelFactoryRepo) Create(ctx context.Context, mf *types.ModelFactory) error {
	return m.db.WithContext(ctx).Create(mf).Error
}

func (m *ModelFactoryRepo) Update(ctx context.Context, mf *types.ModelFactory) error {
	return m.db.WithContext(ctx).Model(mf).Where("id = ?", mf.ID).Updates(mf).Error
}

func (m *ModelFactoryRepo) List(ctx context.Context) ([]*types.ModelFactory, error) {
	var mfs []*types.ModelFactory
	if err := m.db.WithContext(ctx).Model(&types.ModelFactory{}).Scopes(workspaceScope(ctx)).Find(&mfs).Error; err != nil {
		return nil, err
	}
	return mfs, nil
}

func (m *ModelFactoryRepo) Delete(ctx context.Context, id string) error {
	return m.db.WithContext(ctx).Delete(&types.ModelFactory{}, "id = ?", id).Error
}

type ModelRepo struct {
	db *gorm.DB
}

func (m *ModelRepo) CreateMany(ctx context.Context, mfs []*types.Model) error {
	return m.db.WithContext(ctx).Create(mfs).Error
}

func (m *ModelRepo) Create(ctx context.Context, mf *types.Model) error {
	return m.db.WithContext(ctx).Create(mf).Error
}

func (m *ModelRepo) Update(ctx context.Context, mf *types.Model) error {
	return m.db.WithContext(ctx).Model(mf).Where("id = ?", mf.ID).Updates(mf).Error
}

func (m *ModelRepo) List(ctx context.Context, query types.ListModelQuery) ([]*types.Model, error) {
	var ms []*types.Model
	d := m.db.WithContext(ctx).Model(&types.Model{}).Scopes(workspaceScope(ctx))
	if query.FactoryId != "" {
		d = d.Where("factory_id = ?", query.FactoryId)
	}
	if query.Type != "" {
		d = d.Where("type = ?", query.Type)
	}
	if err := d.Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}

func (m *ModelRepo) Delete(ctx context.Context, id string) error {
	return m.db.WithContext(ctx).Delete(&types.Model{}, "id = ?", id).Error
}

func (m *ModelRepo) GetDetail(ctx context.Context, id string) (*types.ModelDetail, error) {
	var md types.ModelDetail
	if err := m.db.WithContext(ctx).
		Model(&types.Model{}).
		Joins("left join model_factories on model_factories.id = models.factory_id").
		Select("models.*, api_key, base_url, model_factories.type as factory_type").
		Where("models.id = ?", id).
		Scopes(workspaceScopeWithTable(ctx, "models")).
		First(&md).Error; err != nil {
		return nil, err
	}
	return &md, nil
}

func NewModelRepo(db *gorm.DB) interfaces.ModelRepo {
	return &ModelRepo{db: db}
}
