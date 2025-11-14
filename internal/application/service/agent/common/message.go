package common

import (
	"beep/internal/models/chat"
	"beep/internal/types"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
)

func ReceiveStream(stream *chat.Stream, msgId, conversationId string, messageChan chan types.Message) (*types.Message, error) {
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
		msg := types.Message{
			ConversationId: conversationId,
			Role:           string(schema.Assistant),
			Content:        chunk.Content,
			BaseEntity:     types.BaseEntity{ID: msgId, CreatedAt: recvTime, UpdatedAt: recvTime},
		}
		sb.WriteString(chunk.Content)
		if len(chunk.ToolCalls) > 0 {
			toolCall := chunk.ToolCalls[0]
			finalMessage.ToolStatus = types.ToolStatusWait
			finalMessage.ToolCall += toolCall.ToolName
			finalMessage.ToolCallParams += toolCall.Arguments
			msg.ToolCall = toolCall.ToolName
			msg.ToolCallParams = toolCall.Arguments
			msg.ToolStatus = types.ToolStatusWait
			if toolCall.ToolCallId != "" {
				msg.ToolCallId = toolCall.ToolCallId
				finalMessage.ToolCallId = toolCall.ToolCallId
			}
		}
		messageChan <- msg
	}
	finalMessage.ID = msgId
	finalMessage.ConversationId = conversationId
	finalMessage.Role = string(schema.Assistant)
	finalMessage.Content = sb.String()
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
