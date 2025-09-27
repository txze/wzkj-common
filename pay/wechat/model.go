package wechat

type PaymentRequest struct {
	OrderId   string `json:"order_id"` //商户订单号
	Expire    string `json:"expire"`
	NotifyUrl string `json:"notify_url"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
}
