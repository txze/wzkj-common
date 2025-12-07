package rabbitmq

type RabbitMQConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	TLS  bool   `json:"tls"`
}

// NewRabbitMQByInternalConf 直接读取nacos的配置，初始化rabbitmq.
func NewRabbitMQWithConfig(serverName string, host, user, pwd string) *RabbitMQ {
	var opts []Option
	opts = append(opts, WithHost(host), WithUser(user), WithPwd(pwd))
	return NewRabbitMQ(
		serverName,
		opts...,
	)
}
