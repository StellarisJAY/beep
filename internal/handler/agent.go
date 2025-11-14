package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	service interfaces.AgentService
}

func NewAgentHandler(service interfaces.AgentService) *AgentHandler {
	return &AgentHandler{
		service: service,
	}
}

func (a *AgentHandler) Create(c *gin.Context) {
	var req types.CreateAgentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	err := a.service.Create(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}
	c.JSON(200, ok())
}

func (a *AgentHandler) List(c *gin.Context) {
	var query types.AgentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	agents, err := a.service.List(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, ok().withData(agents).withTotal(len(agents)))
}

func (a *AgentHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		panic(errors.NewBadRequestError("id不能为空", nil))
	}
	agent, err := a.service.FindById(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, ok().withData(agent))
}
