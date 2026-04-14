package allinpay

// 接口路径常量
const (
	PrepayUrl      = "v1.0/unifiedPay/prepay"       // 主扫支付
	Prepay2Url     = "v1.0/unifiedPay/prepay2"      // JSAPI支付
	FrontPayUrl    = "v1.0/unifiedPay/frontpay"     // 固码支付
	WxSmallUrl     = "v1.0/unifiedPay/wxsmall_pay"  // 微信小程序支付
	MicroUrl       = "v1.0/unifiedPay/micropay"     // 付款码支付
	QueryUrl       = "v1.0/unifiedPay/query"        // 订单查询
	RefundQueryUrl = "v1.0/unifiedPay/refund_query" // 退款查询
	RefundUrl      = "v1.0/unifiedPay/refund"       // 退款
	HistoryUrl     = "history/sm2/order"            // 对账
)

type AllInPayConfig struct {
	AppId      string // 应用ID
	MchntId    string // 商户号
	StoreId    string // 门店号
	ChannelId  string // 渠道号
	PrivateKey string // 私钥
	PublicKey  string // 公钥
	NotifyUrl  string // 支付回调URL
	RefundUrl  string // 退款通知URL
	IsProd     bool   // 是否生产环境
	SignType   string // 签名类型
	PayDomain  string // 支付域名
	HistoryUrl string // 对账接口路径
}

func (a *AllInPayConfig) GetType() string {
	return "allinpay"
}

func (a *AllInPayConfig) GetPayHost() string {
	if a.IsProd {
		return ProdUrl
	}
	return SandboxUrl
}
