package postgres

import "time"

type CommonField struct {
	Id int `gorm:"primary_key;column:id"`

	System    string    `json:"system" gorm:"column:system;INDEX:idx_system"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime;NOT NULL;index:idx_created_at;DEFAULT:CURRENT_TIMESTAMP"`                                // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;NOT NULL;index:idx_updated_at;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0)"` // 最后更新时间
	Status    uint8     `json:"status" gorm:"column:status;NOT NULL;index:idx_status;default:0"`                                                                          // 删除标记，1：normal, 0: delete
}

func GetShortUUID() string {
	db := Client.Master()
	var uuid string
	db.Unscoped().Raw("SELECT uuid_generate_v4()").Row().Scan(&uuid)
	return uuid
}
