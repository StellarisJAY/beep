package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	agentService interfaces.AgentService
}

func NewChatHandler(agentService interfaces.AgentService) *ChatHandler {
	return &ChatHandler{agentService: agentService}
}

func initSSEWriter(c *gin.Context) {
	writer := c.Writer
	writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	_, ok := writer.(http.Flusher)
	if !ok {
		panic(errors.NewBadRequestError("浏览器不支持SSE", nil))
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), types.ContextKeySSEWriter, writer))
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req types.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("请求参数错误", err))
	}
	initSSEWriter(c)
	if req.AgentID != "" || req.Agent != nil {
		if err := h.agentService.Run(c.Request.Context(), types.AgentRunReq{
			AgentId:        req.AgentID,
			Agent:          req.Agent,
			ConversationId: req.ConversationID,
			Query:          req.Query,
		}); err != nil {
			panic(err)
		}
	}
}

func (h *ChatHandler) SignalTool(c *gin.Context) {
	var req types.ToolSignal
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("请求参数错误", err))
	}
	initSSEWriter(c)
	if err := h.agentService.SignalTool(c.Request.Context(), req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}
