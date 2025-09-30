package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ChatService struct {
	conversationRepo interfaces.ConversationRepo
}

func NewChatService(conversationRepo interfaces.ConversationRepo) interfaces.ChatService {
	return &ChatService{conversationRepo: conversationRepo}
}

func (c *ChatService) MessageLoop(ctx context.Context, messageChan chan types.Message, errorChan chan error) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-messageChan:
			if !ok {
				return nil
			}
			if err := c.handleMessage(ctx, msg); err != nil {
				return err
			}
		case err, ok := <-errorChan:
			if !ok {
				return nil
			}
			if err := c.handleError(ctx, err); err != nil {
				return err
			}
		}
	}
}

func (c *ChatService) handleMessage(ctx context.Context, message types.Message) error {
	// 获取SSE Writer
	writer := ctx.Value(types.ContextKeySSEWriter).(gin.ResponseWriter)
	// 发送消息
	if _, err := writer.Write([]byte(fmt.Sprintf("data: %s\n\n", message.Content))); err != nil {
		return errors.NewInternalServerError("发送消息失败", err)
	}
	writer.Flush()
	return nil
}

func (c *ChatService) handleError(ctx context.Context, err error) error {
	return nil
}
