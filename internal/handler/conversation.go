package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	service interfaces.ConversationService
}

func NewConversationHandler(service interfaces.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		service: service,
	}
}

func (h *ConversationHandler) List(c *gin.Context) {
	query := types.ConversationQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	conversations, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, okWithData(conversations))
}

func (h *ConversationHandler) ListMessages(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		panic(errors.NewBadRequestError("id不能为空", nil))
	}
	messages, err := h.service.ListMessages(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, okWithData(messages))
}
