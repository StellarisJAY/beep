package agent

import (
	"beep/internal/application/service/agent/react"
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"

	"github.com/panjf2000/ants/v2"
)

type RunFactory struct {
	modelService         interfaces.ModelService
	memoryService        interfaces.MemoryService
	knowledgeBaseService interfaces.KnowledgeBaseService
	mcpServerService     interfaces.MCPServerService
	conversationRepo     interfaces.ConversationRepo
	messageRepo          interfaces.MessageRepo
	worker               *ants.Pool
}

func NewRunFactory(modelService interfaces.ModelService,
	memoryService interfaces.MemoryService,
	knowledgeBaseService interfaces.KnowledgeBaseService,
	mcpServerService interfaces.MCPServerService,
	conversationRepo interfaces.ConversationRepo,
	messageRepo interfaces.MessageRepo,
	worker *ants.Pool) interfaces.AgentRunFactory {
	return &RunFactory{
		modelService:         modelService,
		memoryService:        memoryService,
		knowledgeBaseService: knowledgeBaseService,
		mcpServerService:     mcpServerService,
		conversationRepo:     conversationRepo,
		messageRepo:          messageRepo,
		worker:               worker,
	}
}

func (f *RunFactory) CreateAgentRun(req types.AgentRunReq) (interfaces.AgentRun, error) {
	if req.Agent.Type == types.AgentTypeReAct {
		return react.NewAgentRun(react.AgentRunParams{
			ModelService:         f.modelService,
			MemoryService:        f.memoryService,
			KnowledgeBaseService: f.knowledgeBaseService,
			McpServerService:     f.mcpServerService,
			Worker:               f.worker,
			ConversationRepo:     f.conversationRepo,
			MessageRepo:          f.messageRepo,
		}, req)
	}
	return nil, errors.NewInternalServerError("不支持的智能体类型", nil)
}
