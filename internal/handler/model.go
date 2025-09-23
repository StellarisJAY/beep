package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	service interfaces.ModelService
}

func NewModelHandler(service interfaces.ModelService) *ModelHandler {
	return &ModelHandler{
		service: service,
	}
}

func (m *ModelHandler) ListModelFactory(c *gin.Context) {
	data, err := m.service.ListFactory(c)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, okWithData(data))
}

func (m *ModelHandler) ListModel(c *gin.Context) {
	var query types.ListModelQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	data, err := m.service.ListModels(c, query)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, okWithData(data))
}

func (m *ModelHandler) CreateModelFactory(c *gin.Context) {
	var req types.CreateModelFactoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := m.service.CreateFactory(c, req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (m *ModelHandler) UpdateFactory(c *gin.Context) {
	var req types.UpdateModelFactoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := m.service.UpdateFactory(c, req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}
