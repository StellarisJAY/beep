package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	service interfaces.DocumentService
}

func NewDocumentHandler(service interfaces.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service: service,
	}
}

func (h *DocumentHandler) CreateFromFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		panic(errors.NewBadRequestError("form file error", err))
	}
	if file == nil {
		panic(errors.NewBadRequestError("文件不能为空", nil))
	}
	kbId, exist := c.GetPostForm("knowledge_base_id")
	if !exist {
		panic(errors.NewBadRequestError("知识库id不能为空", nil))
	}
	id, err := strconv.ParseInt(kbId, 10, 64)
	if err != nil {
		panic(errors.NewBadRequestError("知识库id格式错误", err))
	}
	err = h.service.CreateFromFile(c.Request.Context(), id, file)
	if err != nil {
		panic(err)
	}
	c.JSON(200, ok())
}

func (h *DocumentHandler) List(c *gin.Context) {
	var query types.DocumentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	docs, total, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}
	c.JSON(200, ok().withData(docs).withTotal(total))
}

func (h *DocumentHandler) Rename(c *gin.Context) {
	var req types.RenameDocumentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := h.service.Rename(c.Request.Context(), req); err != nil {
		panic(err)
	}
	c.JSON(200, ok())
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		panic(err)
	}
	c.JSON(200, ok())
}

func (h *DocumentHandler) Download(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	url, err := h.service.Download(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}
	c.JSON(200, ok().withData(url))
}
