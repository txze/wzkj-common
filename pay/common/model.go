package common

type PaymentRequest struct {
	OrderNo   string  `json:"order_no"`   //商户订单号
	Expire    string  `json:"expire"`     //过期时间
	Amount    float64 `json:"amount"`     //支付金额
	Currency  string  `json:"currency"`   //支付货币
	GoodsName string  `json:"goods_name"` //商品名
}

type UnifiedResponse struct {
	Platform    string  `json:"platform"`     // 支付平台
	OrderID     string  `json:"order_id"`     // 商户订单号
	PlatformID  string  `json:"platform_id"`  // 平台订单号
	Amount      float64 `json:"amount"`       // 订单金额
	Status      bool    `json:"status"`       // 支付状态是否成功
	TradeStatus string  `json:"trade_status"` //第三方支付状态
	PaidAmount  float64 `json:"paid_amount"`  // 实付金额
	PaidTime    string  `json:"paid_time"`    // 支付时间
	Message     any     `json:"message"`      // 返回消息
}

type RefundRequest struct {
	OrderNo   string  `json:"order_no"`   //商户订单号
	Amount    float64 `json:"amount"`     //支付金额
	Currency  string  `json:"currency"`   //支付货币
	GoodsName string  `json:"goods_name"` //商品名
}
