package alipay

type AlipayConfig struct {
	AliPayPublicKey string
	Appid           string
	PrivateKey      string
	IsProd          bool
	NotifyUrl       string // 支付回调URL
	RefundUrl       string // 退款通知URL
}

func (a *AlipayConfig) GetType() string {
	return "alipay"
}
