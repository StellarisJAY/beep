package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
)

type AgentService struct {
	repo             interfaces.AgentRepo
	conversationRepo interfaces.ConversationRepo
	agentRunFactory  interfaces.AgentRunFactory
	chatService      interfaces.ChatService
}

func (a *AgentService) Create(ctx context.Context, req types.CreateAgentReq) error {
	agent := &types.Agent{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      *req.Config,
	}
	// 校验智能体配置
	if agent.Type == types.AgentTypeReAct && agent.Config.ReAct == nil {
		return errors.NewBadRequestError("无效的ReACT智能体配置", nil)
	}
	if agent.Type == types.AgentTypeWorkflow && agent.Config.Workflow == nil {
		return errors.NewBadRequestError("无效的Workflow智能配置", nil)
	}
	// 创建智能体
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
	// 查询智能体
	agent, err := a.repo.FindById(ctx, req.Id)
	if err != nil {
		return errors.NewInternalServerError("查询智能体失败", err)
	}
	// 校验智能体配置
	if agent.Type == types.AgentTypeReAct && req.Config.ReAct == nil {
		return errors.NewBadRequestError("无效的ReACT智能配置", nil)
	}
	if agent.Type == types.AgentTypeWorkflow && req.Config.Workflow == nil {
		return errors.NewBadRequestError("无效的Workflow智能配置", nil)
	}
	// 只允许修改智能体名称、描述和配置、配置
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

func (a *AgentService) Run(ctx context.Context, req types.AgentRunReq) error {
	// 如果有指定智能体id，且是普通模式
	// 查询智能体配置
	if req.AgentId != 0 {
		agent, err := a.repo.FindById(ctx, req.AgentId)
		if err != nil {
			return errors.NewInternalServerError("查询智能体失败", err)
		}
		req.Agent = agent
	}

	// 没有指定会话ID，创建一个新会话
	if req.ConversationId == 0 {
		conversation := &types.Conversation{
			AgentId: req.Agent.ID,
		}
		if err := a.conversationRepo.Create(ctx, conversation); err != nil {
			return errors.NewInternalServerError("创建会话失败", err)
		}
		req.ConversationId = conversation.ID
	}

	run, err := a.agentRunFactory.CreateAgentRun(req)
	if err != nil {
		return errors.NewInternalServerError("创建智能体运行失败", err)
	}

	resp, err := run.Run(ctx, req)
	if err != nil {
		return errors.NewInternalServerError("运行智能体失败", err)
	}

	return a.chatService.MessageLoop(ctx, resp.MessageChan, resp.ErrorChan)
}

func NewAgentService(repo interfaces.AgentRepo,
	conversationRepo interfaces.ConversationRepo,
	agentRunFactory interfaces.AgentRunFactory,
	chatService interfaces.ChatService) interfaces.AgentService {
	return &AgentService{repo: repo, conversationRepo: conversationRepo, agentRunFactory: agentRunFactory, chatService: chatService}
}
