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

type PickInfoResponse struct {
	TakeOrderCompanyName        string `json:"TakeOrderCompanyName"`
	TakeOrderCompanyCode        string `json:"TakeOrderCompanyCode"`
	TakeOrderCompanyProvince    string `json:"TakeOrderCompanyProvince"`
	TakeOrderCompanyCity        string `json:"TakeOrderCompanyCity"`
	TakeOrderCompanyArea        string `json:"TakeOrderCompanyArea"`
	TakeOrderCompanyAddress     string `json:"TakeOrderCompanyAddress"`
	TakeOrderOuterPhone         string `json:"TakeOrderOuterPhone"`
	TakeOrderCompanyContactUser string `json:"TakeOrderCompanyContactUser"`
	TakeOrderCourierName        string `json:"TakeOrderCourierName"`
	TakeOrderCourierCode        string `json:"TakeOrderCourierCode"`
	TakeOrderCourierMobile      string `json:"TakeOrderCourierMobile"`
	OrderStatus                 string `json:"OrderStatus"`
	OrderStatusName             string `json:"OrderStatusName"`
	Weight                      string `json:"Weight"`
	PrintCode                   string `json:"PrintCode"`
}
