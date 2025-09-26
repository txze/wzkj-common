package config

type AllInPayConfig struct {
	AppId         string
	CuSID         string
	APIVersion    string
	PrivateKey    string
	PublicKey     string
	NotifyUrl     string
	QueryOrderUrl string
}

func (a *AllInPayConfig) GetType() string {
	return "allinpay"
}
