package shentong

type Config struct {
	IsSandbox    bool //是否开启沙盒 true 是 false 不开启
	AppKey       string
	SecretKey    string
	ResourceCode string
	SourceCode   string
	Customer     Customer //客户信息
}

func (c *Config) GetBaseUrl() string {
	if c.IsSandbox {
		return BaseSandboxUrl
	}
	return BaseUrl
}
