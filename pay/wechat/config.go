package wechat

// mchid：商户ID 或者服务商模式的 sp_mchid
// serialNo：商户API证书的证书序列号
// apiV3Key：APIv3Key，商户平台获取
// privateKey：商户API证书下载后，私钥 apiclient_key.pem 读取后的字符串内容
type WechatConfig struct {
	AppId        string //APPID
	Mchid        string //商户号
	SerialNo     string //商户API证书的证书序列号
	PackageValue string //支付
	ApiV3Key     string //APIv3Key，商户平台获取
	PrivateKey   string //支付私钥
	PublicKey    string //支付公钥
	PublicKeyID  string // 支付公钥ID
}

func (a *WechatConfig) GetType() string {
	return "wechat"
}
