package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MCPServerHandler struct {
	service interfaces.MCPServerService
}

func NewMCPServerHandler(service interfaces.MCPServerService) *MCPServerHandler {
	return &MCPServerHandler{service: service}
}

func (h *MCPServerHandler) Create(c *gin.Context) {
	var req types.CreateMCPServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", nil))
	}
	if err := h.service.Create(c.Request.Context(), req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (h *MCPServerHandler) List(c *gin.Context) {
	result, err := h.service.List(c.Request.Context())
	if err != nil {
		panic(err)
	}
	// TODO 优化，批量获取工具列表
	for _, ms := range result {
		if err := h.service.ListTools(c.Request.Context(), ms); err != nil {
			slog.Debug("MCP服务获取工具列表失败", "url", ms.Url, "err", err)
		}
	}
	c.JSON(http.StatusOK, ok().withData(result))
}

func (h *MCPServerHandler) ListWithoutTools(c *gin.Context) {
	var query types.MCPServerQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(errors.NewBadRequestError("", nil))
	}
	result, err := h.service.ListWithoutTools(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok().withData(result))
}

func (h *MCPServerHandler) Update(c *gin.Context) {
	var req types.UpdateMCPServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", nil))
	}
	if err := h.service.Update(c.Request.Context(), req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (h *MCPServerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}
