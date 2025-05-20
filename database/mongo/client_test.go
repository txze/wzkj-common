package mongo_test

import (
	"context"
	"testing"

	"github.com/txze/wzkj-common/database/mongo"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

type Demo struct {
	Name            string `json:"name" bson:"name"`
	mongo.BaseModel `json:",inline" bson:",inline"`
}

func TestMongoDBOperations(t *testing.T) {
	var username, password, host, port string
	var objectId string
	var err error
	var demo *Demo

	host, port = "127.0.0.1", "27017"
	username, password = "", ""
	mongo.InitMongo(username, password, host, port, "test")

	// 插入一条数据,自动插入时间
	demo = &Demo{
		Name: "tangs",
	}
	_, err = mongo.DB().Collection("test").InsertOne(context.Background(), demo)
	assert.NoError(t, err)

	// 查询数据
	err = mongo.DB().Collection("test").Find(context.Background(), bson.M{
		"name": "tangs",
	}).One(&demo)
	assert.NoError(t, err)
	assert.Equal(t, "tangs", demo.Name)
	assert.Equal(t, demo.CreatedAt.IsZero(), false)
	assert.Equal(t, demo.CreatedAt.IsZero(), false)
	objectId = demo.Id

	// 更新数据， 自动更新updated_at
	err = mongo.DB().Collection("test").UpdateOne(context.Background(),
		bson.M{"_id": objectId},
		bson.M{
			operator.Set: bson.M{
				"name": "Big.tangs",
			},
		},
	)
	assert.NoError(t, err)

	// 查询数据，判断是否更新数据
	err = mongo.DB().Collection("test").Find(context.Background(), bson.M{
		"_id": objectId,
	}).One(&demo)
	assert.NoError(t, err)
	assert.Equal(t, "Big.tangs", demo.Name)
	assert.Equal(t, objectId, demo.Id)

	// 删除数据
	err = mongo.DB().Collection("test").RemoveId(context.Background(), objectId)
	assert.NoError(t, err)

	// 查询数据，判断是否删除成功
	err = mongo.DB().Collection("test").Find(context.Background(), bson.M{
		"_id": objectId,
	}).One(&demo)
	assert.EqualError(t, err, qmgo.ErrNoSuchDocuments.Error())
}
