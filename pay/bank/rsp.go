package bank

type PayResponse struct {
	MchntID    string  `json:"mchntID" mapstructure:"mchntID"`       // 商户ID
	StoreID    string  `json:"storeID" mapstructure:"storeID"`       // 门店ID
	GoodsName  string  `json:"goodsName" mapstructure:"goodsName"`   // 商品名称
	OutOrderNo string  `json:"outOrderNo" mapstructure:"outOrderNo"` // 商户订单号
	TransAmt   float64 `json:"transAmt" mapstructure:"transAmt"`     // 交易金额(元)
	PaySource  string  `json:"paySource" mapstructure:"paySource"`   // 支付来源
	SerProID   string  `json:"serProID" mapstructure:"serProID"`     // 服务产品ID
	SubOpenID  string  `json:"subOpenID" mapstructure:"subOpenID"`   // 用户OpenID
	PayType    string  `json:"payType" mapstructure:"payType"`       // 支付类型
	NotifyURL  string  `json:"notifyURL" mapstructure:"notifyURL"`   // 异步通知地址
	ReturnURL  string  `json:"returnURL" mapstructure:"returnURL"`   // 同步返回地址
	Signature  string  `json:"signature" mapstructure:"signature"`   //签名
}
