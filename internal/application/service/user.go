package service

import (
	"beep/internal/errors"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserServiceImpl struct {
	userRepo          interfaces.UserRepo
	workspaceRepo     interfaces.WorkspaceRepo
	userWorkspaceRepo interfaces.UserWorkspaceRepo
	redis             *redis.Client
}

func NewUserService(userRepo interfaces.UserRepo,
	workspaceRepo interfaces.WorkspaceRepo,
	userWorkspaceRepo interfaces.UserWorkspaceRepo,
	redis *redis.Client) interfaces.UserService {
	return &UserServiceImpl{
		userRepo:          userRepo,
		workspaceRepo:     workspaceRepo,
		userWorkspaceRepo: userWorkspaceRepo,
		redis:             redis,
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

func (u *UserServiceImpl) Login(ctx context.Context, req types.LoginReq) (*types.LoginResp, error) {
	user, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if user == nil {
		return nil, errors.NewError(errors.ErrLoginFailed, "邮箱或密码错误", nil)
	}
	hashPass := sha256.New().Sum([]byte(user.PasswordSalt + req.Password))
	password := hex.EncodeToString(hashPass)
	if password != user.Password {
		return nil, errors.NewError(errors.ErrLoginFailed, "邮箱或密码错误", nil)
	}
	// 登录成功
	// 查询默认工作空间
	workspace, err := u.userWorkspaceRepo.FindUserDefaultWorkspace(ctx, user.ID)
	if err != nil {
		return nil, errors.NewInternalServerError("登录失败", err)
	}
	resp := &types.LoginResp{
		UserInfo:      user,
		WorkspaceInfo: workspace,
		Token:         uuid.New().String(),
		RefreshToken:  uuid.New().String(),
	}
	// 保存登录信息
	loginInfo := &types.LoginInfo{
		UserId:      user.ID,
		WorkspaceId: workspace.ID,
	}
	data, err := json.Marshal(loginInfo)
	if err != nil {
		return nil, errors.NewInternalServerError("登录失败", err)
	}
	// 保存access_token
	if err := u.redis.Set(ctx, "access_token:"+resp.Token, data, time.Hour).Err(); err != nil {
		return nil, errors.NewInternalServerError("登录失败", err)
	}
	// 保存refresh_token
	if err := u.redis.Set(ctx, "refresh_token:"+resp.RefreshToken, user.ID, time.Hour*24).Err(); err != nil {
		return nil, errors.NewInternalServerError("登录失败", err)
	}
	return resp, nil
}
