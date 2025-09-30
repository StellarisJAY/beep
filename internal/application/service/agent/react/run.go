package react

import (
	"beep/internal/models"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"beep/internal/util"
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/panjf2000/ants/v2"
)

type AgentRun struct {
	AgentId        int64
	ConversationId int64
	AgentConfig    types.AgentConfig

	modelService         interfaces.ModelService
	memoryService        interfaces.MemoryService
	knowledgeBaseService interfaces.KnowledgeBaseService
	mcpServerService     interfaces.MCPServerService
	worker               *ants.Pool

	cancelFunc  context.CancelFunc // 取消上下文函数
	iterations  int
	messageChan chan types.Message // 模型回复消息通道
	errChan     chan error         // 错误通道

	chatModel model.BaseChatModel // 聊天模型接口

	tools []*types.MCPToolSet // MCP服务器工具列表
}

type AgentRunParams struct {
	ModelService         interfaces.ModelService
	MemoryService        interfaces.MemoryService
	KnowledgeBaseService interfaces.KnowledgeBaseService
	McpServerService     interfaces.MCPServerService
	Worker               *ants.Pool
}

func NewAgentRun(params AgentRunParams) *AgentRun {
	return &AgentRun{
		modelService:         params.ModelService,
		memoryService:        params.MemoryService,
		knowledgeBaseService: params.KnowledgeBaseService,
		mcpServerService:     params.McpServerService,
		worker:               params.Worker,
	}
}

func (a *AgentRun) Run(ctx context.Context, req types.AgentRunReq) (*types.AgentRunResp, error) {
	a.AgentId = req.Agent.ID
	a.AgentConfig = req.Agent.Config
	a.ConversationId = req.ConversationId
	// 初始化模型回复消息通道和错误通道
	a.messageChan = make(chan types.Message, 128)
	a.errChan = make(chan error, 1)

	// 初始化模型
	chatModelDetail, err := a.modelService.GetModelDetail(ctx, a.AgentConfig.ReAct.ChatModel)
	if err != nil {
		return nil, err
	}
	a.chatModel, err = models.NewChatModel(*chatModelDetail, types.ChatModelOption{})
	if err != nil {
		return nil, err
	}

	// 初始化cancel上下文
	c, cancel := context.WithCancel(ctx)
	a.cancelFunc = cancel

	// 提交任务到worker池
	if err := a.worker.Submit(func() {
		a.reAct(c, req)
	}); err != nil {
		close(a.errChan)
		close(a.messageChan)
		return nil, err
	}

	return &types.AgentRunResp{
		MessageChan: a.messageChan,
		ErrorChan:   a.errChan,
	}, nil
}

func (a *AgentRun) listMcpTools(ctx context.Context) error {
	mcpServerIds := a.AgentConfig.ReAct.McpServers
	mcpServers := make([]*types.MCPToolSet, 0, len(mcpServerIds))
	for _, id := range mcpServerIds {
		mcpServer, err := a.mcpServerService.Get(ctx, id)
		if err != nil {
			slog.Error("获取MCP服务器失败", "id", id, "err", err)
			continue
		}
		err = a.mcpServerService.ListTools(ctx, mcpServer)
		if err != nil {
			slog.Error("获取MCP服务器工具列表失败", "id", id, "err", err)
			continue
		}

		mcpServers = append(mcpServers, &types.MCPToolSet{
			Name:  mcpServer.Name,
			Tools: mcpServer.Tools,
		})
	}
	a.tools = mcpServers
	return nil
}

func (a *AgentRun) reAct(ctx context.Context, req types.AgentRunReq) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("reAct panic", "err", err)
			a.errChan <- err.(error)
		}
	}()
	// 系统提示词
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: a.AgentConfig.ReAct.Prompt,
		},
	}
	// 列出MCP服务器工具列表
	if err := a.listMcpTools(ctx); err != nil {
		slog.Error("获取MCP服务器工具列表失败", "err", err)
	}

	// 添加用户查询
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: req.Query,
	})

	for {
		currentMessage := &types.Message{
			BaseEntity:     types.BaseEntity{ID: util.SnowflakeId()},
			Role:           string(schema.Assistant),
			Content:        "",
			ConversationId: req.ConversationId,
		}
		stream, err := a.chatModel.Stream(ctx, messages)
		if err != nil {
			panic(err)
		}
		for {
			chunk, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				panic(err)
			}
			if chunk.Content != "" {
				currentMessage.Content += chunk.Content
				a.messageChan <- types.Message{
					BaseEntity:     types.BaseEntity{ID: currentMessage.ID},
					Role:           string(schema.Assistant),
					Content:        chunk.Content,
					ToolCall:       chunk.ToolName,
					ConversationId: req.ConversationId,
				}
			}
		}
		messages = append(messages, &schema.Message{
			Role:    schema.Assistant,
			Content: currentMessage.Content,
		})
	}
}

func (a *AgentRun) Cancel(_ context.Context) {
	close(a.messageChan)
	close(a.errChan)
}
