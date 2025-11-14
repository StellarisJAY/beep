package service

import (
	"beep/internal/config"
	"beep/internal/errors"
	"beep/internal/models/chat"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"log/slog"
	"time"
)

type AgentService struct {
	repo              interfaces.AgentRepo
	conversationRepo  interfaces.ConversationRepo
	agentRunFactory   interfaces.AgentRunFactory
	chatService       interfaces.ChatService
	config            *config.Config
	modelService      interfaces.ModelService
	mcpServerRepo     interfaces.MCPServerRepo
	knowledgeBaseRepo interfaces.KnowledgeBaseRepo
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

func (a *AgentService) Detail(ctx context.Context, id string) (*types.Agent, error) {
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

func (a *AgentService) Delete(ctx context.Context, id string) error {
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
	if req.AgentId != "" {
		agent, err := a.repo.FindById(ctx, req.AgentId)
		if err != nil {
			return errors.NewInternalServerError("查询智能体失败", err)
		}
		req.Agent = agent
	}

	// 没有指定会话ID，创建一个新会话
	if req.ConversationId == "" {

		conversation := &types.Conversation{
			AgentId: req.Agent.ID,
			Title:   string([]rune(req.Query)[:min(len(req.Query), 20)]), // 暂时使用查询作为标题
		}
		if err := a.conversationRepo.Create(ctx, conversation); err != nil {
			return errors.NewInternalServerError("创建会话失败", err)
		}
		req.ConversationId = conversation.ID
		// 异步生成会话标题和摘要
		go a.genConversationTitleAndSummary(conversation.ID, req.Query, req.Agent)
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

func (a *AgentService) genConversationTitleAndSummary(conversationId string, query string, agent *types.Agent) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("genConversationTitleAndSummary panic", "err", err)
		}
	}()
	modelId := agent.ChatModelId()
	// 初始化模型
	modelDetail, err := a.modelService.GetModelDetail(context.Background(), modelId)
	if err != nil {
		panic(err)
	}
	chatModel := chat.NewChatModel(*modelDetail)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	messages := []*chat.Message{
		{
			Role:    chat.RoleSystem,
			Content: a.config.ConversationTitlePrompt,
		},
		{
			Role:    chat.RoleUser,
			Content: query,
		},
	}
	// 调用模型
	resp, err := chatModel.Generate(ctx, messages, &chat.Options{})
	if err != nil {
		panic(err)
	}
	// 保存会话标题
	if err := a.conversationRepo.UpdateTitle(ctx, conversationId, resp.Content); err != nil {
		panic(err)
	}
}

func (a *AgentService) SignalTool(ctx context.Context, req types.ToolSignal) error {
	// 没有指定会话ID，创建一个新会话
	if req.ConversationId == "" {
		return errors.NewBadRequestError("会话ID不能为空", nil)
	}

	var agentRunReq types.AgentRunReq

	if req.AgentId != "" {
		agent, err := a.repo.FindById(ctx, req.AgentId)
		if err != nil {
			return errors.NewInternalServerError("查询智能体失败", err)
		}
		agentRunReq.AgentId = req.AgentId
		agentRunReq.Agent = agent
		agentRunReq.ConversationId = req.ConversationId
	}

	run, err := a.agentRunFactory.CreateAgentRun(agentRunReq)
	if err != nil {
		return errors.NewInternalServerError("创建智能体运行失败", err)
	}

	resp, err := run.SignalTool(ctx, req)
	if err != nil {
		return errors.NewInternalServerError("调用工具失败", err)
	}
	return a.chatService.MessageLoop(ctx, resp.MessageChan, resp.ErrorChan)
}

func (a *AgentService) FindById(ctx context.Context, id string) (*types.AgentDetail, error) {
	agent, err := a.repo.FindById(ctx, id)
	if err != nil {
		return nil, errors.NewInternalServerError("查询智能体失败", err)
	}

	detail := &types.AgentDetail{
		Agent: *agent,
	}

	if len(agent.McpServerIds()) > 0 {
		// 查询关联的MCP服务器
		mcpServers, err := a.mcpServerRepo.ListWithoutTools(ctx, types.MCPServerQuery{
			Ids: agent.McpServerIds(),
		})
		if err != nil {
			return nil, errors.NewInternalServerError("查询MCP服务器失败", err)
		}
		detail.McpServers = mcpServers
	}

	if len(agent.KnowledgeBaseIds()) > 0 {
		// 查询关联的知识库
		knowledgeBases, _, err := a.knowledgeBaseRepo.List(ctx, types.KnowledgeBaseQuery{
			Ids: agent.KnowledgeBaseIds(),
		})
		if err != nil {
			return nil, errors.NewInternalServerError("查询知识库失败", err)
		}
		detail.KnowledgeBases = knowledgeBases
	}
	return detail, nil
}

func NewAgentService(repo interfaces.AgentRepo,
	conversationRepo interfaces.ConversationRepo,
	agentRunFactory interfaces.AgentRunFactory,
	chatService interfaces.ChatService,
	config *config.Config,
	modelService interfaces.ModelService,
	mcpServerRepo interfaces.MCPServerRepo,
	knowledgeBaseRepo interfaces.KnowledgeBaseRepo,
) interfaces.AgentService {
	return &AgentService{
		repo:              repo,
		conversationRepo:  conversationRepo,
		agentRunFactory:   agentRunFactory,
		chatService:       chatService,
		config:            config,
		modelService:      modelService,
		mcpServerRepo:     mcpServerRepo,
		knowledgeBaseRepo: knowledgeBaseRepo,
	}
}
