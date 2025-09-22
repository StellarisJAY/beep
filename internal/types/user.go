package types

import "time"

type User struct {
	BaseEntity
	Name          string     `json:"name" gorm:"not null;type:varchar(64);'"`
	Email         string     `json:"email" gorm:"unique;not null;type:varchar(255);"`
	Password      string     `json:"password" gorm:"not null;type:varchar(255)"`
	PasswordSalt  string     `json:"password_salt" gorm:"not null;type:varchar(255)"`
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

type LoginInfo struct {
	User
	WorkspaceId int64 `json:"workspace_id;string"`
}
