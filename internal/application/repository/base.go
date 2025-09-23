package repository

import "gorm.io/gorm"

func pageScope(paged bool, pageNum, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if paged {
			return db.Offset((pageNum - 1) * pageSize).Limit(pageSize)
		}
		return db
	}
}
