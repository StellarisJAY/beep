package types

import (
	"beep/internal/util"
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        int64          `json:"id,string" gorm:"primary_key'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		b.ID = util.SnowflakeId()
	}
	return nil
}
