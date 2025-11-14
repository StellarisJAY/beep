package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type AgentRepo struct {
	db *gorm.DB
}

func (a *AgentRepo) Create(ctx context.Context, agent *types.Agent) error {
	return a.db.WithContext(ctx).Create(agent).Error
}

func (a *AgentRepo) FindById(ctx context.Context, id string) (*types.Agent, error) {
	var agent *types.Agent
	err := a.db.WithContext(ctx).Scopes(workspaceScope(ctx)).Where("id = ?", id).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (a *AgentRepo) Update(ctx context.Context, agent *types.Agent) error {
	return a.db.WithContext(ctx).Model(agent).Scopes(workspaceScope(ctx)).Where("id=?", agent.ID).Updates(agent).Error
}

func (a *AgentRepo) Delete(ctx context.Context, id string) error {
	return a.db.WithContext(ctx).Scopes(workspaceScope(ctx)).Delete(&types.Agent{}, "id=?", id).Error
}

func (a *AgentRepo) List(ctx context.Context, query types.AgentQuery) ([]*types.Agent, error) {
	var agents []*types.Agent
	d := a.db.WithContext(ctx).Model(&types.Agent{}).
		Scopes(workspaceScope(ctx)).
		Select("id, name, type, description, created_at, updated_at, create_by, workspace_id")
	if query.Name != "" {
		d = d.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type != "" {
		d = d.Where("type = ?", query.Type)
	}
	if query.CreateByMe {
		d = d.Where("create_by = ?", ctx.Value(types.UserIdContextKey).(int64))
	}
	if err := d.Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func NewAgentRepo(db *gorm.DB) interfaces.AgentRepo {
	return &AgentRepo{db: db}
}
