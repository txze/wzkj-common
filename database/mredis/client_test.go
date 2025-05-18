package mredis_test

import (
	"testing"

	"wzkj-common/database/mredis"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	var uri = "127.0.0.1:6379"
	var password = "Jj5uXPvo"
	var db = 0

	var err error
	var value int
	var reply interface{}

	mredis.Init(uri, password, db)

	// 设置key test-redit
	_, err = mredis.Master().Do("set", "test-redis", 1)
	assert.NoError(t, err)

	// 获取key
	reply, err = mredis.Master().Do("get", "test-redis")
	assert.NoError(t, err)
	value, err = redis.Int(reply, err)
	assert.NoError(t, err)
	assert.Equal(t, 1, value)
}
