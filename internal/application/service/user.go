package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type UserServiceImpl struct {
	userRepo          interfaces.UserRepo
	workspaceRepo     interfaces.WorkspaceRepo
	userWorkspaceRepo interfaces.UserWorkspaceRepo
}

func NewUserService(userRepo interfaces.UserRepo, workspaceRepo interfaces.WorkspaceRepo, userWorkspaceRepo interfaces.UserWorkspaceRepo) interfaces.UserService {
	return &UserServiceImpl{
		userRepo:          userRepo,
		workspaceRepo:     workspaceRepo,
		userWorkspaceRepo: userWorkspaceRepo,
	}
}

func (u *UserServiceImpl) Register(ctx context.Context, req types.RegisterReq) error {
	exist, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if exist != nil {
		return errors.NewConflictError("邮箱已经注册", nil)
	}
	// 密码加密
	salt := uuid.New().String()[:16]
	hashPass := sha256.New().Sum([]byte(salt + req.Password))
	user := &types.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     hex.EncodeToString(hashPass),
		PasswordSalt: salt,
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return errors.NewInternalServerError("注册失败", err)
	}
	// 创建用户默认工作空间
	workspace := &types.Workspace{
		Name:        fmt.Sprintf("%s的工作空间", user.Name),
		Description: fmt.Sprintf("%s的初始工作空间", user.Name),
	}
	if err := u.workspaceRepo.Create(ctx, workspace); err != nil {
		return errors.NewInternalServerError("创建工作空间失败", err)
	}
	// 创建用户默认工作空间关联
	uw := &types.UserWorkspace{
		UserId:      user.ID,
		WorkspaceId: workspace.ID,
		Role:        types.WorkspaceRoleOwner,
	}
	if err := u.userWorkspaceRepo.Create(ctx, uw); err != nil {
		return errors.NewInternalServerError("创建用户工作空间失败", err)
	}
	return nil
}

func (u *UserServiceImpl) Login(ctx context.Context, req types.LoginReq) (*types.User, error) {
	panic("not implemented")
}
