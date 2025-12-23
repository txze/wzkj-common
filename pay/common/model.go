package common

import "time"

type PaymentRequest struct {
	OrderNo   string `json:"order_no"`   //商户订单号
	Expire    string `json:"expire"`     //过期时间
	Amount    int    `json:"amount"`     //支付金额
	Currency  string `json:"currency"`   //支付货币
	GoodsName string `json:"goods_name"` //商品名
	Params    string `json:"params"`     //回传参数
}

type UnifiedResponse struct {
	Platform       string    `json:"platform"`        // 支付平台
	OrderID        string    `json:"order_id"`        // 商户订单号
	PlatformID     string    `json:"platform_id"`     // 平台订单号
	Amount         int       `json:"amount"`          // 订单金额
	Status         bool      `json:"status"`          // 支付状态是否成功
	TradeStatus    string    `json:"trade_status"`    //第三方支付状态
	PaidAmount     int       `json:"paid_amount"`     // 实付金额
	PaidTime       time.Time `json:"paid_time"`       // 支付时间
	Message        any       `json:"message"`         // 返回消息
	Params         string    `json:"params"`          //回传参数
	DiscountAmount int       `json:"discount_amount"` // 优惠金额
}

type RefundRequest struct {
	RefundNo  string `json:"refund_no"`  // 退款流水号
	OrderNo   string `json:"order_no"`   // 商户订单号
	Amount    int    `json:"amount"`     // 支付金额
	Total     int    `json:"total"`      // 愿支付金额
	Currency  string `json:"currency"`   // 支付货币
	GoodsName string `json:"goods_name"` // 商品名
}

type RefundResponse struct {
	UserReceivedAccount  string `json:"user_received_account"`  // 退款入账账户
	SuccessTime          string `json:"success_time"`           // 退款成功时间
	CreateTime           string `json:"create_time"`            // 退款创建时间
	OriginalRefundStatus string `json:"original_refund_status"` // 退款状态
	RefundStatus         bool   `json:"refund_status"`          // 退款状态
	Message              string `json:"message"`                //错误信息
	RefundAmount         int    `json:"amount"`
	Data                 any    `json:"data"` // 退款数据
}

type RefundOrderResponse struct {
	OutRefundNo         string `json:"out_refund_no"`         // 商户退款单号
	TransactionId       string `json:"transaction_id"`        // 支付系统生成的订单号
	OutTradeNo          string `json:"out_trade_no"`          // 商户系统内部订单号
	Channel             string `json:"channel"`               // 退款渠道
	UserReceivedAccount string `json:"user_received_account"` // 退款入账账户
	SuccessTime         string `json:"success_time"`          // 退款成功时间
	CreateTime          string `json:"create_time"`           // 退款创建时间
	Status              string `json:"status"`                // 退款状态
	Total               int    `json:"total"`                 // 订单总金额，单位为分
	Refund              int    `json:"refund"`                // 退款标价金额，单位为分，可以做部分退款
	PayerTotal          int    `json:"payer_total"`           // 用户支付金额，单位为分
	PayerRefund         int    `json:"payer_refund"`          // 用户退款金额，不包含所有优惠券金额
	IsSuccess           bool   `json:"is_success"`            //是否成功
	RefundInfo          any    `json:"refund_info"`           //退款信息
}
