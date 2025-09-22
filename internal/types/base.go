package types

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uint64         `json:"id,string" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	return nil
}
