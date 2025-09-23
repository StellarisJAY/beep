package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	service interfaces.WorkspaceService
}

func NewWorkspaceHandler(service interfaces.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{
		service: service,
	}
}

func (w *WorkspaceHandler) List(c *gin.Context) {
	userId := c.GetInt64("user_id")
	workspaces, err := w.service.ListByUserId(c.Request.Context(), userId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, okWithData(workspaces))
}

func (w *WorkspaceHandler) ListMember(c *gin.Context) {
	workspaceId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	members, err := w.service.ListMembers(c.Request.Context(), workspaceId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, okWithData(members))
}

func (w *WorkspaceHandler) InviteMember(c *gin.Context) {
	var req types.InviteWorkspaceMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := w.service.InviteMember(c.Request.Context(), req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (w *WorkspaceHandler) SwitchWorkspace(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := w.service.SwitchWorkspace(c.Request.Context(), id); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (w *WorkspaceHandler) SetRole(c *gin.Context) {
	var req types.SetWorkspaceRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := w.service.SetRole(c.Request.Context(), req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}
