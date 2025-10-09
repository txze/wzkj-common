package alipay

type AlipayConfig struct {
	AliPayPublicKey         string
	Appid                   string
	PrivateKey              string
	IsProd                  bool
	NotifyUrl               string // 支付回调URL
	RefundUrl               string // 退款通知URL
	AppCertContent          []byte // 应用公钥证书文件内容
	AliPayRootCertContent   []byte //支付宝根证书文件内容
	AliPayPublicCertContent []byte //支付宝公钥证书文件内容

}

func (a *AlipayConfig) GetType() string {
	return "alipay"
}
