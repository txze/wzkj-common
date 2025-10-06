package alipay

type AlipayConfig struct {
	AliPayPublicKey string
	Appid           string
	PrivateKey      string
	IsProd          bool
}

func (a *AlipayConfig) GetType() string {
	return "alipay"
}
