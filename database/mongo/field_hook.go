package mongo

import (
	"time"

	"github.com/qiniu/qmgo/field"
)

type BaseModel struct {
	Id        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Status    int       `json:"status" bson:"status"`
}

func (bm *BaseModel) CustomFields() field.CustomFieldsBuilder {
	return field.NewCustom().
		SetCreateAt("CreatedAt").
		SetUpdateAt("UpdatedAt").
		SetId("Id")
}

// func (bm *BaseModel) BeforeInsert() error {
// 	bm.Id = qmgo.NewObjectID().Hex()
// 	bm.CreatedAt = util.Now()
// 	bm.CreatedAt = util.Now()
// 	bm.Status = define.StatusNormal
// 	return nil
// }

// func (bm *BaseModel) AfterInsert() error {
// 	return nil
// }

// func (bm *BaseModel) BeforeUpdate() error {

// }

// func (bm *BaseModel) BeforeQuery() error {

// }

// func (bm *BaseModel) AfterQuery() error {

// }

// func (bm *BaseModel) BeforeRemove() error {

// }

// func (bm *BaseModel) AfterRemove() error {

// }

// func (bm *BaseModel) BeforeUpsert() error {

// }

// func (bm *BaseModel) AfterUpsert() error {

// }
