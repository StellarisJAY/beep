package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
)

type AgentService struct {
	repo interfaces.AgentRepo
}

func (a *AgentService) Create(ctx context.Context, req types.CreateAgentReq) error {
	agent := &types.Agent{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      *req.Config,
	}
	if agent.Type == types.AgentTypeReAct && agent.Config.ReAct == nil {
		return errors.NewBadRequestError("无效的ReACT智能体配置", nil)
	}
	if agent.Type == types.AgentTypeWorkflow && agent.Config.Workflow == nil {
		return errors.NewBadRequestError("无效的Workflow智能配置", nil)
	}
	if err := a.repo.Create(ctx, agent); err != nil {
		return errors.NewInternalServerError("创建智能体失败", err)
	}
	return nil
}

func (a *AgentService) Detail(ctx context.Context, id int64) (*types.Agent, error) {
	agent, err := a.repo.FindById(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerError("查询智能体失败", err)
	}
	return agent, nil
}

func (a *AgentService) Update(ctx context.Context, req types.UpdateAgentReq) error {
	agent, err := a.repo.FindById(ctx, req.Id)
	if err != nil {
		return errors.NewInternalServerError("查询智能体失败", err)
	}
	if agent.Type == types.AgentTypeReAct && req.Config.ReAct == nil {
		return errors.NewBadRequestError("无效的ReACT智能配置", nil)
	}
	if agent.Type == types.AgentTypeWorkflow && req.Config.Workflow == nil {
		return errors.NewBadRequestError("无效的Workflow智能配置", nil)
	}
	agent.Name = req.Name
	agent.Description = req.Description
	agent.Config = *req.Config
	if err := a.repo.Update(ctx, agent); err != nil {
		return errors.NewInternalServerError("更新智能体失败", err)
	}
	return nil
}

func (a *AgentService) Delete(ctx context.Context, id int64) error {
	if err := a.repo.Delete(ctx, id); err != nil {
		return errors.NewInternalServerError("删除智能体失败", err)
	}
	return nil
}

func (a *AgentService) List(ctx context.Context, query types.AgentQuery) ([]*types.Agent, error) {
	agents, err := a.repo.List(ctx, query)
	if err != nil {
		return nil, errors.NewInternalServerError("查询智能体列表失败", err)
	}
	return agents, nil
}

func (a *AgentService) Run(ctx context.Context, id int64, query string) error {
	//TODO implement me
	panic("implement me")
}

func NewAgentService(repo interfaces.AgentRepo) interfaces.AgentService {
	return &AgentService{repo: repo}
}
