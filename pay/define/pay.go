package define

// 支付平台类型常量
const (
	PlatformAlipay = "alipay" // 支付宝
	PlatformWechat = "wechat" // 微信
)

// 分账模式常量
const (
	RoyaltyModeAsync = "async" // 异步分账
	RoyaltyModeSync  = "sync"  // 同步分账
)

// 支付状态
const (
	StatusSuccess = "SUCCESS" // 成功
	StatusFail    = "FAIL"    // 失败
	StatusPending = "PENDING" // 支付中
	StatusCancel  = "CANCEL"  // 已取消
	StatusRefund  = "REFUND"  // 转入退款
	StatusClose   = "CLOSE"   // 已关闭 交易关闭，不可退款 全部退款成功
)
