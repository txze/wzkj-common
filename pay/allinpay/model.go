package allinpay

import (
	"github.com/hzxiao/goutil"
)

// 交易状态常量
const (
	TradeStatusSuccess = "SUCCESS" // 支付成功
	TradeStatusFailed  = "FAILED"  // 支付失败
	TradeStatusPending = "PENDING" // 处理中
	TradeStatusClosed  = "CLOSED"  // 已关闭
)

// 签名类型常量
const (
	SignTypeSM2 = "SM2"
	SignTypeRSA = "RSA"
)

// 支付来源常量
const (
	PaySourceDefault = "0" // 默认支付来源
	PaySourceWap     = "1" // WAP支付
)

// 支付方式常量
const (
	PayTypeWechat = "1" // 微信支付
	PayTypeAlipay = "2" // 支付宝
	PayTypeSmall  = "3" // 微信小程序
	PayTypeApp    = "4" // 应用支付
)

const (
	// 沙盒环境
	SandboxUrl = "https://zjk.numrd.com/aps/buybal-api/"
	// 生产环境
	ProdUrl = "https://s.zjkccb.com/aps/buybal-api/"
)

// 支付请求参数
type PayRequest struct {
	GoodsName  string `json:"goodsName"`  // 商品名称
	OutOrderNo string `json:"outOrderNo"` // 商户订单号
	TransAmt   int    `json:"transAmt"`   // 交易金额（分）
	PaySource  string `json:"paySource"`  // 支付来源
	PayType    string `json:"payType"`    // 支付方式
	AppId      string `json:"appId"`      // 应用ID（微信小程序支付）
	OpenId     string `json:"openId"`     // 用户ID（微信小程序支付）
	ReturnUrl  string `json:"returnUrl"`  // 返回URL（固码支付）
}

// 退款请求参数
type RefundRequest struct {
	OutOrderNo   string `json:"outOrderNo"`   // 商户订单号
	OrderNo      string `json:"orderNo"`      // 平台订单号
	OutRefundNo  string `json:"outRefundNo"`  // 退款订单号
	RefundAmount string `json:"refundAmount"` // 退款金额（分）
}

// 查询请求参数
type QueryRequest struct {
	OutOrderNo string `json:"outOrderNo"` // 商户订单号
}

// 统一响应结构
type ApiResponse struct {
	RespCode string     `json:"respCode"`  // 响应码
	RespDesc string     `json:"respDesc"`  // 响应信息
	Data     goutil.Map `json:"data"`      // 响应数据
	Sign     string     `json:"signature"` // 签名
}

// 回调通知结构体
type NotifyRequest struct {
	MchntId     string `json:"mchntId"`     // 商户号
	StoreId     string `json:"storeId"`     // 门店号
	ChannelId   string `json:"channelId"`   // 服务商号
	OutOrderNo  string `json:"outOrderNo"`  // 外部订单号
	TransAmt    int    `json:"transAmt"`    // 总金额
	OrderNo     string `json:"orderNo"`     // 订单号
	PayTime     string `json:"payTime"`     // 支付完成时间
	PaySource   int    `json:"paySource"`   // 支付通道
	PayType     int    `json:"payType"`     // 支付方式
	TradeStatus int    `json:"tradeStatus"` // 交易状态
	SignType    string `json:"signType"`    // 签名算法
	Signature   string `json:"signature"`   // 签名
	ContractId  string `json:"contractId"`  // 签约协议号（银联无感支付）
}
