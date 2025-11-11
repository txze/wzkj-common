package shentong

type CreateOrderResponse struct {
	OrderNo       string `mapstructure:"orderNo"`
	WaybillNo     string `mapstructure:"waybillNo"`
	BigWord       string `mapstructure:"bigWord"`
	PackagePlace  string `mapstructure:"packagePlace"`
	SourceOrderId string `mapstructure:"sourceOrderId"`
	SafeNo        string `mapstructure:"safeNo"`
	NewBlockCode  string `mapstructure:"newBlockCode"`
}
