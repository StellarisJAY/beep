package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KnowledgeBaseHandler struct {
	service interfaces.KnowledgeBaseService
}

func NewKnowledgeBaseHandler(service interfaces.KnowledgeBaseService) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{
		service: service,
	}
}

func (k *KnowledgeBaseHandler) Create(c *gin.Context) {
	var req types.CreateKnowledgeBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := k.service.Create(c, req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (k *KnowledgeBaseHandler) List(c *gin.Context) {
	var query types.KnowledgeBaseQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	res, total, err := k.service.List(c, query)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, okWithData(res).withTotal(total))
}

func (k *KnowledgeBaseHandler) Update(c *gin.Context) {
	var req types.UpdateKnowledgeBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := k.service.Update(c, req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (k *KnowledgeBaseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		panic(errors.NewBadRequestError("id不能为空", nil))
	}
	if err := k.service.Delete(c, id); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}
