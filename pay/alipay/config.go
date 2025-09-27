package alipay

type AlipayConfig struct {
	AliPayPublicKey string
}

func (a *AlipayConfig) GetType() string {
	return "alipay"
}
