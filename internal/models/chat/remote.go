package chat

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

type remoteAPIModel struct {
	client    openai.Client
	modelName string
}

// Generate 聊天文本生成
func (r *remoteAPIModel) Generate(ctx context.Context, messages []*Message, options *Options) (*Message, error) {
	params := r.makeCompletionParams(messages, options)
	response, err := r.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}
	choice := response.Choices[0]
	assistantMessage := &Message{
		Role:          string(choice.Message.Role),
		Content:       choice.Message.Content,
		FunctionCalls: convertOpenaiToolCalls(choice.Message.ToolCalls),
	}
	return assistantMessage, nil
}

// Stream 聊天流式输出
func (r *remoteAPIModel) Stream(ctx context.Context, messages []*Message, options *Options) (*Stream, error) {
	params := r.makeCompletionParams(messages, options)
	stream := r.client.Chat.Completions.NewStreaming(ctx, params)
	outputStream := newStream()
	// 启动一个goroutine处理流式输出
	go func() {
		defer outputStream.Close()
		defer stream.Close()

		for stream.Next() {
			chunk := stream.Current()
			// 错误输出到errChan
			if err := stream.Err(); err != nil {
				outputStream.errChan <- err
				return
			}
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta
			message := &Message{
				Role:          delta.Role,
				Content:       delta.Content,
				FunctionCalls: convertOpenaiToolCallsStream(delta.ToolCalls),
			}
			outputStream.messageChan <- message
		}
	}()
	return outputStream, nil
}

// getChatTools 将MCP服务和内部工具转换成OpenAI的Tool定义
func getChatTools(options *Options) []openai.ChatCompletionToolUnionParam {
	tools := make([]openai.ChatCompletionToolUnionParam, 0, len(options.McpServers))
	for _, mcpServer := range options.McpServers {
		for _, mcpTool := range mcpServer.Tools {
			tool := &openai.ChatCompletionFunctionToolParam{
				Type: "function",
				Function: shared.FunctionDefinitionParam{
					// 每个工具的名字为 mcp:mcp服务名:工具名
					Name:        fmt.Sprintf("mcp:%s:%s", mcpServer.Name, mcpTool.Name),
					Strict:      param.NewOpt(true),
					Description: param.NewOpt(mcpTool.Description),
					Parameters:  nil, // TODO 处理参数
				},
			}
			tools = append(tools, openai.ChatCompletionToolUnionParam{OfFunction: tool})
		}
	}
	return tools
}

// convertMessages 将Message转换为OpenAI的ChatCompletionMessage
func convertMessages(messages []*Message) []openai.ChatCompletionMessageParamUnion {
	convertedMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, message := range messages {
		var convertedMessage openai.ChatCompletionMessageParamUnion
		switch message.Role {
		case "user":
			convertedMessage.OfUser = &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{OfString: param.NewOpt(message.Content)},
			}
		case "assistant":
			convertedMessage.OfAssistant = &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{OfString: param.NewOpt(message.Content)},
			}
		case "system":
			convertedMessage.OfSystem = &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{OfString: param.NewOpt(message.Content)},
			}
		}
		convertedMessages = append(convertedMessages, convertedMessage)
	}
	return convertedMessages
}

// makeCompletionParams 构建OpenAI的ChatCompletionNewParams
func (r *remoteAPIModel) makeCompletionParams(messages []*Message, options *Options) openai.ChatCompletionNewParams {
	// 转换输入消息转换成OpenAI格式
	inputMessages := convertMessages(messages)
	params := openai.ChatCompletionNewParams{
		Model: r.modelName, // 模型名称
		// 模型调用参数
		Temperature:     param.NewOpt(options.Temperature),
		TopP:            param.NewOpt(options.TopP),
		PresencePenalty: param.NewOpt(options.FrequencyPenalty),
		// 输入消息
		Messages: inputMessages,
	}
	// 获取工具定义
	tools := getChatTools(options)
	if len(tools) > 0 {
		params.Tools = tools
	}
	return params
}

// convertOpenaiToolCallsStream 将OpenAI的流式工具调用转换为ToolCall
func convertOpenaiToolCallsStream(toolCalls []openai.ChatCompletionChunkChoiceDeltaToolCall) []*ToolCall {
	convertedToolCalls := make([]*ToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		convertedToolCalls = append(convertedToolCalls, &ToolCall{
			ToolName:  toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		})
	}
	return convertedToolCalls
}

// convertOpenaiToolCalls 将OpenAI的工具调用转换为ToolCall
func convertOpenaiToolCalls(toolCalls []openai.ChatCompletionMessageToolCallUnion) []*ToolCall {
	convertedToolCalls := make([]*ToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		convertedToolCalls = append(convertedToolCalls, &ToolCall{
			ToolName:  toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		})
	}
	return convertedToolCalls
}
