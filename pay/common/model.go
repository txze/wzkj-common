package common

type PaymentRequest struct {
	OrderNo   string `json:"order_no"`   //商户订单号
	Expire    string `json:"expire"`     //过期时间
	Amount    int    `json:"amount"`     //支付金额
	Currency  string `json:"currency"`   //支付货币
	GoodsName string `json:"goods_name"` //商品名
	Params    string `json:"params"`     //回传参数
}

type UnifiedResponse struct {
	Platform    string `json:"platform"`     // 支付平台
	OrderID     string `json:"order_id"`     // 商户订单号
	PlatformID  string `json:"platform_id"`  // 平台订单号
	Amount      int    `json:"amount"`       // 订单金额
	Status      bool   `json:"status"`       // 支付状态是否成功
	TradeStatus string `json:"trade_status"` //第三方支付状态
	PaidAmount  int    `json:"paid_amount"`  // 实付金额
	PaidTime    string `json:"paid_time"`    // 支付时间
	Message     any    `json:"message"`      // 返回消息
	Params      string `json:"params"`       //回传参数
}

type RefundRequest struct {
	RefundNo  string `json:"refund_no"`  // 退款流水号
	OrderNo   string `json:"order_no"`   // 商户订单号
	Amount    int    `json:"amount"`     // 支付金额
	Currency  string `json:"currency"`   // 支付货币
	GoodsName string `json:"goods_name"` // 商品名
}

type RefundResponse struct {
	UserReceivedAccount string `json:"user_received_account"` // 退款入账账户
	SuccessTime         string `json:"success_time"`          // 退款成功时间
	CreateTime          string `json:"create_time"`           // 退款创建时间
}
