package goredis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	master *redis.Client
	slave  *redis.Client
}

var defaultClient = &Client{}

func NewRedisClient(uri string, password string, db int) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
		DB:       db,

		// ========== 连接池核心配置 ==========
		PoolSize:        50,               // 根据实际并发量调整
		MinIdleConns:    5,                // 最小空闲连接数
		MaxIdleConns:    20,               // 最大空闲连接数
		PoolTimeout:     3 * time.Second,  // 获取连接超时
		ConnMaxIdleTime: 10 * time.Minute, // 空闲连接超时
		ConnMaxLifetime: 20 * time.Minute, // 连接最大存活时间

		// ========== 超时配置 ==========
		DialTimeout:  5 * time.Second, // 建立连接超时
		ReadTimeout:  3 * time.Second, // 读超时
		WriteTimeout: 3 * time.Second, // 写超时

		// ========== 重试配置 ==========
		MaxRetries:      3,                      // 最大重试次数
		MinRetryBackoff: 100 * time.Millisecond, // 最小重试间隔
		MaxRetryBackoff: 1 * time.Second,        // 最大重试间隔

	})

	defaultClient.master = rdb
	defaultClient.slave = rdb
}

func Master() *redis.Client {
	return defaultClient.master
}

func Salve() *redis.Client {
	return defaultClient.slave
}
