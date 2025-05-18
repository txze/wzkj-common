package mysql

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type CommonField struct {
	Id        int                   `json:"id" gorm:"primary_key;column:id"`
	System    string                `json:"system" gorm:"column:system;INDEX:idx_system"`
	CreatedAt time.Time             `json:"created_at" gorm:"column:created_at;type:datetime;NOT NULL;index:idx_created_at;DEFAULT:CURRENT_TIMESTAMP"`                                // 创建时间
	UpdatedAt time.Time             `json:"updated_at" gorm:"column:updated_at;type:datetime;NOT NULL;index:idx_updated_at;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0)"` // 最后更新时间
	Deleted   soft_delete.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;softDelete:flag;default:0"`
}

func GetShortUUID() string {
	db := Client.Master()
	var uuid string
	db.Unscoped().Raw("SELECT UUID_SHORT()").Row().Scan(&uuid)
	return uuid
}
