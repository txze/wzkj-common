package common

import (
	"time"

	"github.com/go-pay/gopay"

	"github.com/txze/wzkj-common/pay/define"
)

// TradeRoyaltyRateQueryRequestInterface 分账查询请求接口
type TradeRoyaltyRateQueryRequestInterface interface {
	ToMap() gopay.BodyMap
}

// TradeRoyaltyRateQueryResponse 统一的分账查询响应结构
type TradeRoyaltyRateQueryResponse struct {
	// 平台类型：alipay 或 wechat
	Platform            string `json:"platform"`
	TransactionId       string `json:"transaction_id"`        // 交易订单号
	OutOrderNo          string `json:"out_order_no"`          // 分账订单号
	TransactionSettleNo string `json:"transaction_settle_no"` // 分账结算号
	Code                string `json:"code"`
	Msg                 string `json:"msg"`

	// 微信分账响应
	// 只有微信存在
	SubMchid string `json:"sub_mchid"` //电商平台二级商户号
	// 原始响应数据，用于调试和扩展
	RawData interface{} `json:"raw_data,omitempty"`
}

type PaymentRequest struct {
	OrderNo     string `json:"order_no"`     //商户订单号
	Expire      string `json:"expire"`       //过期时间
	Amount      int    `json:"amount"`       //支付金额
	Currency    string `json:"currency"`     //支付货币
	GoodsName   string `json:"goods_name"`   //商品名
	Params      string `json:"params"`       //回传参数
	ProductCode string `json:"product_code"` // 产品码
	// 是否为支付通/服务商模式
	IsServiceMode     bool                `json:"is_service_mode"` // 是否为支付通/服务商模式
	SettleDetailInfos []map[string]string `json:"settle_info"`     // 分账资金信息
	SubMerchantID     string              `json:"sub_merchant_id"` // 子商户ID
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

// RoyaltyDetail 分账明细
type RoyaltyDetail struct {
	OperationType  string `json:"operation_type"`    // 分账操作类型：transfer(分账)、transfer_refund(退分账)、replenish(补差)、replenish_refund(退补差)
	Amount         int    `json:"amount"`            // 分账金额，单位为分
	State          string `json:"state"`             // 分账状态：SUCCESS(成功)、FAIL(失败)、PROCESSING(处理中)
	ExecuteDt      string `json:"execute_dt"`        // 分账执行时间
	TransOut       string `json:"trans_out"`         // 分账转出账号
	TransOutType   string `json:"trans_out_type"`    // 分账转出账号类型
	TransOutOpenId string `json:"trans_out_open_id"` // 分账转出方的openId
	TransIn        string `json:"trans_in"`          // 分账转入账号
	TransInType    string `json:"trans_in_type"`     // 分账转入账号类型
	TransInOpenId  string `json:"trans_in_open_id"`  // 分账转入方的openId
	DetailId       string `json:"detail_id"`         // 分账明细单号
	ErrorCode      string `json:"error_code"`        // 分账失败错误码
	ErrorDesc      string `json:"error_desc"`        // 分账错误描述信息
}

// SettleNotificationResponse 分账通知响应
type SettleNotificationResponse struct {
	Platform            string          `json:"platform"`              // 平台类型：alipay 或 wechat
	OutRequestNo        string          `json:"out_request_no"`        // 商户分账请求号
	MsgType             string          `json:"msg_type"`              // 消息类型：AUTO_SETTLE_FINISH(超期自动分账解冻)、ASYNC_SETTLE_RESULT(异步分账结果)
	TradeNo             string          `json:"trade_no"`              // 交易号
	RoyaltyFinishAmount int             `json:"royalty_finish_amount"` // 分账完结金额，单位为分
	SettleNo            string          `json:"settle_no"`             // 支付宝分账受理单号
	OperationDt         string          `json:"operation_dt"`          // 分账受理时间
	OperationFinishDt   string          `json:"operation_finish_dt"`   // 业务执行完成时间
	RoyaltyDetailList   []RoyaltyDetail `json:"royalty_detail_list"`   // 分账明细
	RawData             interface{}     `json:"raw_data"`              // 原始响应数据，用于调试和扩展
}

// SettleConfirmRequestInterface 结算确认请求接口
type SettleConfirmRequestInterface interface {
	GetPlatform() string
}

// SettleConfirmResponse 结算确认响应
type SettleConfirmResponse struct {
	Platform string      `json:"platform"`           // 平台类型：alipay 或 wechat
	Code     string      `json:"code"`               // 响应码
	Msg      string      `json:"msg"`                // 响应信息
	SubCode  string      `json:"sub_code"`           // 子响应码
	SubMsg   string      `json:"sub_msg"`            // 子响应信息
	RawData  interface{} `json:"raw_data,omitempty"` // 原始响应数据，用于调试和扩展
}

// 微信/支付宝状态转换
func ConvertPaymentStatus(platform, originalStatus string) string {
	switch platform {
	case define.PlatformAlipay:
		return convertAlipayStatus(originalStatus)
	case define.PlatformWechat:
		return convertWechatStatus(originalStatus)
	default:
		return define.StatusFail
	}
}

// 支付宝状态转换
func convertAlipayStatus(status string) string {
	switch status {
	case "TRADE_SUCCESS":
		return define.StatusSuccess
	case "TRADE_FINISHED":
		return define.StatusSuccess
	case "WAIT_BUYER_PAY":
		return define.StatusPending
	case "TRADE_CLOSED":
		return define.StatusClose
	default:
		return define.StatusFail
	}
}

// 微信状态转换
func convertWechatStatus(status string) string {
	switch status {
	case "SUCCESS":
		return define.StatusSuccess
	case "USERPAYING":
		return define.StatusPending
	case "NOTPAY":
		return define.StatusPending
	case "REFUND":
		return define.StatusRefund
	case "CLOSED":
		return define.StatusClose
	case "REVOKED":
		return define.StatusCancel
	case "PAYERROR":
		return define.StatusFail
	default:
		return define.StatusFail
	}
}

// 分账状态转换
func ConvertRoyaltyStatus(status string) string {
	switch status {
	case "SUCCESS":
		return define.StatusSuccess
	case "PROCESSING":
		return define.StatusPending
	case "FAIL":
		return define.StatusFail
	default:
		return define.StatusFail
	}
}
