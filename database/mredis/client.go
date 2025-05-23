package mredis

import (
	"github.com/gomodule/redigo/redis"
)

type Client struct {
	master *redis.Pool
	slave  *redis.Pool
}

var defalutClient = &Client{}

func Init(uri string, password string, db int) {
	var options []redis.DialOption
	if password != "" {
		options = append(options, redis.DialPassword(password))
	}
	options = append(options, redis.DialDatabase(db))

	var pool = &redis.Pool{ //实例化一个连接池
		MaxIdle: 16, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			conn, err := redis.Dial("tcp", uri, options...)
			if err != nil {
				return nil, err
			}

			return conn, err
		},
	}

	defalutClient.master = pool
	defalutClient.slave = pool
}

func Master() redis.Conn {
	return defalutClient.master.Get()
}

func Salve() redis.Conn {
	return defalutClient.slave.Get()
}
