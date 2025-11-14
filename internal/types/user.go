package types

import "time"

const (
	UserIdContextKey      = "user_id"
	WorkspaceIdContextKey = "workspace_id"
	AccessTokenContextKey = "access_token"
)

// User 用户
type User struct {
	BaseEntity
	Name          string     `json:"name" gorm:"not null;type:varchar(64);'"`         // 用户名
	Email         string     `json:"email" gorm:"unique;not null;type:varchar(255);"` // 邮箱
	Password      string     `json:"-" gorm:"not null;type:varchar(255)"`             // 密码
	PasswordSalt  string     `json:"-" gorm:"not null;type:varchar(255)"`             // 密码盐
	LastLoginTime *time.Time `json:"last_login_time"`                                 // 最后登录时间
	LastLoginIp   string     `json:"last_login_ip;"`                                  // 最后登录IP
}

func (User) TableName() string {
	return "users"
}

// RegisterReq 注册请求
type RegisterReq struct {
	Name     string `json:"name" binding:"required"`     // 用户名
	Email    string `json:"email" binding:"required"`    // 邮箱
	Password string `json:"password" binding:"required"` // 密码
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`    // 邮箱
	Password string `json:"password" binding:"required"` // 密码
}

type LoginResp struct {
	UserInfo      *User      `json:"user_info"`      // 用户信息
	WorkspaceInfo *Workspace `json:"workspace_info"` // 工作空间信息
	Token         string     `json:"token"`          // 访问令牌
	RefreshToken  string     `json:"refresh_token"`  // 刷新令牌
}

type LoginInfo struct {
	UserId      string `json:"user_id"`
	WorkspaceId string `json:"workspace_id"`
}
