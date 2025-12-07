package rabbitmq

import (
	"errors"

	"glprivate.quanyougame.net/backend/core-sdk-go/nacos"
)

type RabbitMQConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	TLS  bool   `json:"tls"`
}

func loadConfig() (*RabbitMQConf, error) {
	var conf map[string]map[string]*RabbitMQConf
	if err := nacos.ScanInfraConf(&conf); err != nil {
		return nil, err
	}
	mqConf, ok := conf["mq"]
	if !ok {
		return nil, errors.New("mq config not found")
	}
	rmqConf, ok := mqConf["rabbitmq"]
	if !ok {
		return nil, errors.New("rabbitmq config not found")
	}
	return rmqConf, nil
}

// NewRabbitMQByInternalConf 直接读取nacos的配置，初始化rabbitmq.
func NewRabbitMQByInternalConf(serverName string, opts ...Option) *RabbitMQ {
	conf, err := loadConfig()
	if err != nil {
		panic("failed to load rabbitmq config: " + err.Error())
	}
	opts = append(opts, WithHost(conf.Host), WithUser(conf.User), WithPwd(conf.Pwd))
	if !conf.TLS {
		opts = append(opts, WithDisableTLS(true))
	}
	return NewRabbitMQ(
		serverName,
		opts...,
	)
}
