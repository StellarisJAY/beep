package handler

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService interfaces.UserService
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) Register(c *gin.Context) {
	var req types.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	if err := u.userService.Register(c.Request.Context(), req); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, ok())
}

func (u *UserHandler) Login(c *gin.Context) {
	var req types.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errors.NewBadRequestError("", err))
	}
	resp, err := u.userService.Login(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}
	c.Header("Access-Token", resp.Token)
	c.Header("Refresh-Token", resp.RefreshToken)
	c.JSON(http.StatusOK, okWithData(resp))
}

func (u *UserHandler) GetLoginInfo(c *gin.Context) {
	loginInfo, err := u.userService.GetLoginInfo(c.Request.Context())
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, okWithData(loginInfo))
}
