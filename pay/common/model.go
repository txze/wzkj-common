package common

type PaymentRequest struct {
	OrderNo   string `json:"order_no"`   //商户订单号
	Expire    string `json:"expire"`     //过期时间
	NotifyUrl string `json:"notify_url"` //支付后的回调
	Amount    int64  `json:"amount"`     //支付金额
	Currency  string `json:"currency"`   //支付货币
	GoodsName string `json:"goods_name"` //商品名
}
