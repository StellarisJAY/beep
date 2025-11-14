package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type MCPServerRepo struct {
	db *gorm.DB
}

func NewMCPServerRepo(db *gorm.DB) interfaces.MCPServerRepo {
	return &MCPServerRepo{db: db}
}

func (m *MCPServerRepo) Create(ctx context.Context, ms *types.MCPServer) error {
	return m.db.WithContext(ctx).Create(ms).Error
}

func (m *MCPServerRepo) Update(ctx context.Context, ms *types.MCPServer) error {
	return m.db.Model(ms).WithContext(ctx).Where("id=?", ms.ID).Updates(ms).Error
}

func (m *MCPServerRepo) Delete(ctx context.Context, id string) error {
	return m.db.WithContext(ctx).Delete(&types.MCPServer{}, "id = ?", id).Error
}

func (m *MCPServerRepo) List(ctx context.Context) ([]*types.MCPServer, error) {
	var ms []*types.MCPServer
	if err := m.db.WithContext(ctx).Model(&types.MCPServer{}).Scopes(workspaceScope(ctx)).Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}

func (m *MCPServerRepo) Get(ctx context.Context, id string) (*types.MCPServer, error) {
	var ms *types.MCPServer
	if err := m.db.WithContext(ctx).Model(&types.MCPServer{}).Scopes(workspaceScope(ctx)).First(&ms, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return ms, nil
}

func (m *MCPServerRepo) ListWithoutTools(ctx context.Context, query types.MCPServerQuery) ([]*types.MCPServer, error) {
	var ms []*types.MCPServer
	d := m.db.WithContext(ctx).Model(&types.MCPServer{}).Scopes(workspaceScope(ctx))
	if query.Name != "" {
		d = d.Where("name like ?", "%"+query.Name+"%")
	}
	if query.Url != "" {
		d = d.Where("url like ?", "%"+query.Url+"%")
	}
	if len(query.Ids) > 0 {
		d = d.Where("id in ?", query.Ids)
	}
	if err := d.Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}
