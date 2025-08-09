package bank

import "github.com/txze/wzkj-common/pay/common"

type ConfigBank struct {
	Common     common.PaymentCommonConfig
	PrivateKey string // 私钥
	PublicKey  string // 公钥
	MchntID    string // 商户ID
}
