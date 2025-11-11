package react

import (
	"beep/internal/application/service/agent/common"
	"beep/internal/models/chat"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"beep/internal/util"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	conversationRepo     interfaces.ConversationRepo
	messageRepo          interfaces.MessageRepo
	worker               *ants.Pool

	cancelFunc  context.CancelFunc // 取消上下文函数
	iterations  int
	messageChan chan types.Message // 模型回复消息通道
	errChan     chan error         // 错误通道

	chatModel chat.BaseModel // 聊天模型接口

	tools []*types.MCPServer // MCP服务器工具列表
}

type AgentRunParams struct {
	ModelService         interfaces.ModelService
	MemoryService        interfaces.MemoryService
	KnowledgeBaseService interfaces.KnowledgeBaseService
	McpServerService     interfaces.MCPServerService
	ConversationRepo     interfaces.ConversationRepo
	MessageRepo          interfaces.MessageRepo
	Worker               *ants.Pool
}

func NewAgentRun(params AgentRunParams, req types.AgentRunReq) (*AgentRun, error) {
	// 初始化模型
	chatModelDetail, err := params.ModelService.GetModelDetail(context.Background(), req.Agent.Config.ReAct.ChatModel)
	if err != nil {
		return nil, err
	}
	run := &AgentRun{
		modelService:         params.ModelService,
		memoryService:        params.MemoryService,
		knowledgeBaseService: params.KnowledgeBaseService,
		mcpServerService:     params.McpServerService,
		worker:               params.Worker,
		conversationRepo:     params.ConversationRepo,
		messageRepo:          params.MessageRepo,
		AgentId:              req.AgentId,
		AgentConfig:          req.Agent.Config,
		ConversationId:       req.ConversationId,
		messageChan:          make(chan types.Message, 128),
		errChan:              make(chan error, 1),
		chatModel:            chat.NewChatModel(*chatModelDetail),
	}

	if err := run.listMcpTools(context.TODO()); err != nil {
		return nil, err
	}
	return run, nil
}

func (a *AgentRun) Run(ctx context.Context, req types.AgentRunReq) (*types.AgentRunResp, error) {
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

func (a *AgentRun) SignalTool(ctx context.Context, signal types.ToolSignal) (*types.AgentRunResp, error) {
	// 查询上一次模型回复的工具调用消息
	isTool := true
	lastMessages, err := a.messageRepo.Search(ctx, types.MessageQuery{
		ConversationId: signal.ConversationId,
		Role:           chat.RoleAssistant,
		ToolCall:       &isTool,
	})
	if err != nil {
		return nil, err
	}
	// 检查是否有工具调用消息
	if len(lastMessages) == 0 {
		return nil, errors.New("没有工具调用消息")
	}
	lastMessage := lastMessages[len(lastMessages)-1]
	// 调用工具
	parts := strings.SplitN(lastMessage.ToolCall, ":", 3)
	cmd, toolSet, toolName := parts[0], parts[1], parts[2]
	result, err := a.callTool(ctx, cmd, toolSet, toolName, lastMessage.ToolCallParams)
	if err != nil {
		return nil, err
	}

	// 初始化cancel上下文
	c, cancel := context.WithCancel(ctx)
	a.cancelFunc = cancel
	// 提交任务到worker池
	if err := a.worker.Submit(func() {
		a.toolReply(c, result, lastMessage.ToolCallId)
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

func (a *AgentRun) toolReply(ctx context.Context, toolResult string, toolCallId string) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("reAct panic", "err", err)
			a.errChan <- err.(error)
		}
		close(a.messageChan)
		close(a.errChan)
	}()
	// 让模型根据工具调用结果回复
	// 系统提示词
	messages := []*chat.Message{
		{
			Role:    chat.RoleSystem,
			Content: a.AgentConfig.ReAct.Prompt,
		},
	}
	// 读取记忆
	messages = append(messages, a.loadMemory(ctx, int(a.AgentConfig.ReAct.MemoryOption.MemoryWindowSize+2))...)
	// 工具调用结果消息
	toolMessage := &chat.Message{
		Role:       chat.RoleTool,
		Content:    toolResult,
		ToolCallId: toolCallId,
	}
	messages = append(messages, toolMessage)
	_ = a.messageRepo.Create(ctx, &types.Message{
		Role:           chat.RoleTool,
		Content:        toolResult,
		ConversationId: a.ConversationId,
		ToolCallId:     toolCallId,
	})

	a.sendAndReceive(ctx, messages)
}

// listMcpTools 将agent可用的MCP服务转成工具调用的格式
func (a *AgentRun) listMcpTools(ctx context.Context) error {
	mcpServerIds := a.AgentConfig.ReAct.McpServers
	mcpServers := make([]*types.MCPServer, 0, len(mcpServerIds))
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
		mcpServers = append(mcpServers, mcpServer)
	}
	a.tools = mcpServers
	return nil
}

// callTool 进行工具调用，返回工具结果的json字符串。cmd可能是mcp或local，分别代表MCP远程调用和本地工具
func (a *AgentRun) callTool(ctx context.Context, cmd, toolSet, toolName string, paramsJson string) (string, error) {
	if cmd == "mcp" {
		idx := slices.IndexFunc(a.tools, func(s *types.MCPServer) bool {
			return s.Name == toolSet
		})
		if idx == -1 {
			return "", errors.New("MCP服务器不存在")
		}
		tool := a.tools[idx]
		params := make(map[string]string)
		_ = json.Unmarshal([]byte(paramsJson), &params)
		result, err := a.mcpServerService.Call(ctx, tool.ID, &mcp.CallToolParams{
			Name:      toolName,
			Arguments: params,
		})
		if err != nil {
			return "", err
		}
		slog.Info("MCP服务器调用成功", "tool_set", toolSet, "tool_name", toolName, "arguments", paramsJson, "result", result.Content)
		if len(result.Content) == 0 {
			return "", nil
		}
		resultJSON, err := result.Content[0].MarshalJSON()
		if err != nil {
			return "", err
		}
		return string(resultJSON), nil
	}
	// TODO 本地工具调用
	return "", nil
}

func (a *AgentRun) reAct(ctx context.Context, req types.AgentRunReq) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("reAct panic", "err", err)
			a.errChan <- err.(error)
		}
		close(a.messageChan)
		close(a.errChan)
	}()

	// 系统提示词
	messages := []*chat.Message{
		{
			Role:    chat.RoleSystem,
			Content: a.AgentConfig.ReAct.Prompt,
		},
	}
	// 读取记忆
	messages = append(messages, a.loadMemory(ctx, int(a.AgentConfig.ReAct.MemoryOption.MemoryWindowSize))...)

	// 添加用户查询
	messages = append(messages, &chat.Message{
		Role:    chat.RoleUser,
		Content: req.Query,
	})

	// 回送用户消息
	userMessage := types.Message{
		Role:           chat.RoleUser,
		Content:        req.Query,
		ConversationId: req.ConversationId,
	}
	if err := a.messageRepo.Create(ctx, &userMessage); err != nil {
		panic(err)
	}
	a.messageChan <- userMessage

	a.sendAndReceive(ctx, messages)

}

func (a *AgentRun) sendAndReceive(ctx context.Context, messages []*chat.Message) {
	// 发送消息，接受stream
	currentMessage := &types.Message{
		BaseEntity:     types.BaseEntity{ID: util.SnowflakeId()},
		Role:           chat.RoleAssistant,
		ConversationId: a.ConversationId,
	}

	stream, err := a.chatModel.Stream(ctx, messages, &chat.Options{
		McpServers: a.tools,
	})
	if err != nil {
		panic(err)
	}

	finalMessage, err := common.ReceiveStream(stream, currentMessage.ID, a.ConversationId, a.messageChan)
	if err != nil {
		panic(err)
	}
	// 保存消息记录
	if err := a.messageRepo.Create(ctx, finalMessage); err != nil {
		panic(err)
	}
}

func (a *AgentRun) Cancel(_ context.Context) {
	close(a.messageChan)
	close(a.errChan)
}

func (a *AgentRun) loadMemory(ctx context.Context, windowSize int) []*chat.Message {
	// 读取记忆
	memoryOption := a.AgentConfig.ReAct.MemoryOption
	if memoryOption.MemoryControl == types.MemoryControlStatic {
		if memoryOption.EnableShortTermMemory {
			memoryMessages, err := common.GetStaticMemory(ctx, a.memoryService, a.ConversationId, windowSize)
			if err != nil {
				slog.Error("获取短期记忆失败",
					"conversation_id", a.ConversationId,
					"window_size", int(memoryOption.MemoryWindowSize),
					"err", err)
			}
			return memoryMessages
		}
	}
	return []*chat.Message{}
}
