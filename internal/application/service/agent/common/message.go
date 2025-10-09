package common

import (
	"beep/internal/types"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
)

func ReceiveStream(stream *schema.StreamReader[*schema.Message], msgId, conversationId int64, messageChan chan types.Message) (*types.Message, error) {
	recvTime := time.Now()
	sb := new(strings.Builder)
	finalMessage := &types.Message{}
	for {
		chunk, err := stream.Recv()
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
			finalMessage.ToolCall = toolCall.Function.Name
			finalMessage.ToolCallParams = toolCall.Function.Arguments
			messageChan <- *finalMessage
			return finalMessage, nil
		}
		if chunk.Content != "" {
			sb.WriteString(chunk.Content)
			messageChan <- types.Message{
				BaseEntity:     types.BaseEntity{ID: msgId, CreatedAt: recvTime, UpdatedAt: recvTime},
				Role:           string(schema.Assistant),
				Content:        chunk.Content,
				ToolCall:       chunk.ToolName,
				ConversationId: conversationId,
			}
		}
	}
	finalMessage.ID = msgId
	finalMessage.ConversationId = conversationId
	finalMessage.Role = string(schema.Assistant)
	finalMessage.Content = sb.String()
	messageChan <- *finalMessage
	return finalMessage, nil
}

func ConvertMessageToSchemaMessage(msg *types.Message) *schema.Message {
	return &schema.Message{
		Role:     schema.RoleType(msg.Role),
		Content:  msg.Content,
		ToolName: msg.ToolCall,
	}
}
