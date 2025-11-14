package types

import (
	"beep/internal/util"
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        string         `json:"id" gorm:"primary_key;type:varchar(36);"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = util.UUID()
	}
	return nil
}

type BaseQuery struct {
	Paged    bool `form:"paged"`
	PageNum  int  `form:"page_num"`
	PageSize int  `form:"page_size"`
}
