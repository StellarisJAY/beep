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
	worker               *ants.Pool
}

func NewRunFactory(modelService interfaces.ModelService,
	memoryService interfaces.MemoryService,
	knowledgeBaseService interfaces.KnowledgeBaseService,
	mcpServerService interfaces.MCPServerService,
	worker *ants.Pool) interfaces.AgentRunFactory {
	return &RunFactory{
		modelService:         modelService,
		memoryService:        memoryService,
		knowledgeBaseService: knowledgeBaseService,
		mcpServerService:     mcpServerService,
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
		}), nil
	}
	return nil, errors.NewInternalServerError("不支持的智能体类型", nil)
}
