package config

// PaymentConfig 支付配置接口
type PaymentConfig interface {
	GetType() string
}

// PaymentCommonConfig  公共配置
type PaymentCommonConfig struct {
	NotifyURL     string // 异步回调地址
	SyncReturnURL string // 同步返回地址
}
