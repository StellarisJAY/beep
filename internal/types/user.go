package types

import "time"

const (
	UserIdContextKey      = "user_id"
	WorkspaceIdContextKey = "workspace_id"
	AccessTokenContextKey = "access_token"
)

type User struct {
	BaseEntity
	Name          string     `json:"name" gorm:"not null;type:varchar(64);'"`
	Email         string     `json:"email" gorm:"unique;not null;type:varchar(255);"`
	Password      string     `json:"-" gorm:"not null;type:varchar(255)"`
	PasswordSalt  string     `json:"-" gorm:"not null;type:varchar(255)"`
	LastLoginTime *time.Time `json:"last_login_time"`
	LastLoginIp   string     `json:"last_login_ip;type:varchar(16)"`
}

func (User) TableName() string {
	return "users"
}

type RegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	UserInfo      *User      `json:"user_info"`
	WorkspaceInfo *Workspace `json:"workspace_info"`
	Token         string     `json:"token"`
	RefreshToken  string     `json:"refresh_token"`
}

type LoginInfo struct {
	UserId      int64 `json:"user_id;string"`
	WorkspaceId int64 `json:"workspace_id;string"`
}
