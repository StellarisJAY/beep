package types

import "time"

type User struct {
	BaseEntity
	Name          string     `json:"name" gorm:"not null"`
	Email         string     `json:"email" gorm:"unique;not null"`
	Password      string     `json:"password" gorm:"not null"`
	PasswordSalt  string     `json:"password_salt" gorm:"not null"`
	LastLoginTime *time.Time `json:"last_login_time"`
	LastLoginIp   string     `json:"last_login_ip"`
}

func (User) TableName() string {
	return "users"
}
