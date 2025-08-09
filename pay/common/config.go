package common

// PaymentCommonConfig  公共配置
type PaymentCommonConfig struct {
	NotifyURL     string // 异步回调地址
	SyncReturnURL string // 同步返回地址
}
