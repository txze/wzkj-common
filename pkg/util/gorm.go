package util

import (
	"gorm.io/gorm"
)

// Where 根据条件是否拼接where查询语句
func Where(db *gorm.DB, cond bool, query interface{}, args ...interface{}) *gorm.DB {
	if !cond {
		return db
	}
	return db.Where(query, args...)
}

// Order 排序拼接
func Order(db *gorm.DB, value interface{}) *gorm.DB {
	return db.Order(value)
}

// OffsetAndLimit 拼接offset和limit
func OffsetAndLimit(db *gorm.DB, offset, limit int) *gorm.DB {
	if offset != 0 {
		db = db.Offset(offset)
	}
	if limit != 0 {
		db = db.Limit(limit)
	}
	return db
}
