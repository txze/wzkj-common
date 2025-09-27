package wechat

type WechatConfig struct {
	AppId        string
	Mchid        string
	PackageValue string
	ApiV3Key     string
}

func (a *WechatConfig) GetType() string {
	return "wechat"
}
