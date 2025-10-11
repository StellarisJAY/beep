package chat

import (
	"beep/internal/types"
	"context"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"io"
)

// ToolCall 函数调用
type ToolCall struct {
	Type      string `json:"type"`      // 工具类型，function,mcp
	ToolName  string `json:"tool_name"` // 工具名称
	Arguments string `json:"arguments"` // 工具参数，JSON字符串
}

// Message 聊天消息
type Message struct {
	Role      string      `json:"role"`       // 消息角色，user/assistant/system
	Content   string      `json:"content"`    // 消息内容
	ToolCalls []*ToolCall `json:"tool_calls"` // 工具调用列表
}

// Options 模型调用参数
type Options struct {
	Reasoning        bool    `json:"reasoning"`         // 是否开启推理模式
	MaxTokens        int     `json:"max_tokens"`        // 最大输出token数
	Temperature      float64 `json:"temperature"`       // 温度参数，控制输出的随机性
	TopP             float64 `json:"top_p"`             // Top-p采样，控制输出的多样性
	Stop             string  `json:"stop"`              // 停止生成的token
	FrequencyPenalty float64 `json:"frequency_penalty"` // 频率惩罚参数，控制输出的重复度

	McpServers        []*types.MCPServer `json:"-"`                   // 模型可用的MCP服务器列表
	ParallelToolCalls bool               `json:"parallel_tool_calls"` // 是否并行调用工具
}

// BaseModel 聊天模型基础接口，适配不同的模型API，如OpenAI、Ollama等
type BaseModel interface {
	// Generate 生成回复
	Generate(ctx context.Context, messages []*Message, options *Options) (*Message, error)
	// Stream 流式生成回复
	Stream(ctx context.Context, messages []*Message, options *Options) (*Stream, error)
}

// NewChatModel 创建一个聊天模型
func NewChatModel(modelDetail types.ModelDetail) BaseModel {
	switch modelDetail.FactoryType {
	case types.FactoryDashscope, types.FactoryOpenAI, types.FactoryOpenAICompatible:
		chatModel := &remoteAPIModel{
			client:    openai.NewClient(option.WithBaseURL(modelDetail.BaseUrl), option.WithAPIKey(modelDetail.ApiKeyDecrypted)),
			modelName: modelDetail.Name,
		}
		return chatModel
	case types.FactoryOllama:
		panic("not implemented")
	default:
		panic("not implemented")
	}
}

type Stream struct {
	messageChan chan *Message
	errChan     chan error
}

func newStream() *Stream {
	return &Stream{
		messageChan: make(chan *Message, 64),
		errChan:     make(chan error, 1),
	}
}

func (s *Stream) Next() (*Message, error) {
	select {
	case message, ok := <-s.messageChan:
		if !ok {
			return nil, io.EOF
		}
		return message, nil
	case err, ok := <-s.errChan:
		if !ok {
			return nil, io.EOF
		}
		return nil, err
	}
}

func (s *Stream) Close() {
	close(s.messageChan)
	close(s.errChan)
}
