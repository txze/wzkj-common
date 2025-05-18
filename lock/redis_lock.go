package lock

import (
	"fmt"
	"wzkj-common/database/mredis"
	"wzkj-common/pkg/ierr"

	"github.com/gomodule/redigo/redis"
)

// StartConcurrentLimit 并发锁
func StartConcurrentLimit(prefix string, value string) error {
	// 处理并发问题
	var redisKey = fmt.Sprintf("%s_%s", prefix, value)
	reply, err := mredis.Master().Do("Incr", redisKey)
	if err != nil {
		return ierr.NewIError(ierr.RedisOperatorFail, err.Error())
	}
	// 设置超时时间
	_, _ = mredis.Master().Do("expire", redisKey, 20)

	keyNum, err := redis.Int(reply, err)
	if err != nil {
		return ierr.NewIError(ierr.RedisOperatorFail, err.Error())
	}
	if keyNum > 1 {
		return ierr.NewIError(ierr.NetworkBusyError, "网络繁忙，请稍后再试")
	}

	return nil
}

// CancelConcurrentLimit 删除并发锁
func CancelConcurrentLimit(prefix string, value string) error {
	var err error
	var redisKey = fmt.Sprintf("%s_%s", prefix, value)
	_, err = mredis.Master().Do("DEL", redisKey)
	if err != nil {
		return ierr.NewIError(ierr.RedisOperatorFail, err.Error())
	}

	return nil
}
