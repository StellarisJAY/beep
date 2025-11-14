package interfaces

import (
	"beep/internal/types"
	"context"
)

// UserRepo 用户数据库
type UserRepo interface {
	// Create 创建用户
	Create(ctx context.Context, user *types.User) error
	// Update 更新用户
	Update(ctx context.Context, user *types.User) error
	// Delete 删除用户
	Delete(ctx context.Context, user *types.User) error
	// FindById 根据ID查找用户
	FindById(ctx context.Context, userId string) (*types.User, error)
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*types.User, error)
	// CheckPassword 检查密码是否正确
	CheckPassword(ctx context.Context, email string, password string) (*types.User, error)
}

// UserService 用户服务
type UserService interface {
	// Register 注册用户
	Register(ctx context.Context, req types.RegisterReq) error
	// Login 用户登录
	Login(ctx context.Context, req types.LoginReq) (*types.LoginResp, error)
	// GetLoginInfo 获取用户登录信息
	GetLoginInfo(ctx context.Context) (*types.LoginResp, error)
}
