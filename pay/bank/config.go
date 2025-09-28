package bank

import (
	"github.com/txze/wzkj-common/pay/config"
)

type ConfigBank struct {
	Common     config.PaymentCommonConfig
	PrivateKey string // 私钥
	PublicKey  string // 公钥
	MchntID    string // 商户ID
}
