package common

type UnifiedResponse struct {
	Platform   string `json:"platform"`    // 支付平台
	OrderID    string `json:"order_id"`    // 商户订单号
	PlatformID string `json:"platform_id"` // 平台订单号
	Amount     int    `json:"amount"`      // 订单金额
	Status     string `json:"status"`      // 支付状态
	PaidAmount int    `json:"paid_amount"` // 实付金额
	PaidTime   string `json:"paid_time"`   // 支付时间
	Message    any    `json:"message"`     // 返回消息
}
