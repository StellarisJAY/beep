package common

import (
	"beep/internal/models/chat"
	"beep/internal/types"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
)

func ReceiveStream(stream *chat.Stream, msgId, conversationId int64, messageChan chan types.Message) (*types.Message, error) {
	recvTime := time.Now()
	sb := new(strings.Builder)
	finalMessage := &types.Message{}
	for {
		chunk, err := stream.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(chunk.ToolCalls) > 0 {
			toolCall := chunk.ToolCalls[0]
			finalMessage.ID = msgId
			finalMessage.ConversationId = conversationId
			finalMessage.Role = string(schema.Assistant)
			finalMessage.Content = chunk.Content
			finalMessage.ToolCall += toolCall.ToolName
			finalMessage.ToolCallParams += toolCall.Arguments
			if toolCall.ToolCallId != "" {
				finalMessage.ToolCallId = toolCall.ToolCallId
			}
		}
		if chunk.Content != "" {
			sb.WriteString(chunk.Content)
			messageChan <- types.Message{
				BaseEntity:     types.BaseEntity{ID: msgId, CreatedAt: recvTime, UpdatedAt: recvTime},
				Role:           string(schema.Assistant),
				Content:        chunk.Content,
				ConversationId: conversationId,
			}
		}
	}
	finalMessage.ID = msgId
	finalMessage.ConversationId = conversationId
	finalMessage.Role = string(schema.Assistant)
	if finalMessage.ToolCall != "" {
		parts := strings.SplitN(finalMessage.ToolCall, ":", 3)
		if len(parts) != 3 {
			return nil, errors.New("invalid tool call")
		}
		cmd, toolSet, toolName := parts[0], parts[1], parts[2]
		finalMessage.Content = fmt.Sprintf("请求调用工具 %s:%s，调用方式:%s", toolSet, toolName, cmd)
	} else {
		finalMessage.Content = sb.String()
	}
	messageChan <- *finalMessage
	return finalMessage, nil
}

func convertMessage(message *types.Message) *chat.Message {
	msg := &chat.Message{
		Role:       message.Role,
		Content:    message.Content,
		ToolCallId: message.ToolCallId,
	}
	if message.ToolCallId != "" {
		msg.ToolCallId = message.ToolCallId
		msg.ToolCalls = []*chat.ToolCall{
			{
				ToolCallId: message.ToolCallId,
				ToolName:   message.ToolCall,
				Arguments:  message.ToolCallParams,
			},
		}
	}
	return msg
}
