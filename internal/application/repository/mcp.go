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

func (m *MCPServerRepo) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&types.MCPServer{}, "id = ?", id).Error
}

func (m *MCPServerRepo) List(ctx context.Context) ([]*types.MCPServer, error) {
	var ms []*types.MCPServer
	if err := m.db.WithContext(ctx).Model(&types.MCPServer{}).Scopes(workspaceScope(ctx)).Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}

func (m *MCPServerRepo) Get(ctx context.Context, id int64) (*types.MCPServer, error) {
	var ms *types.MCPServer
	if err := m.db.WithContext(ctx).Model(&types.MCPServer{}).Scopes(workspaceScope(ctx)).First(&ms, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return ms, nil
}
