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

type QuerySendServiceDetailResponse struct {
	AvailableServiceItemList []struct {
		FeeModel FeeModel `mapstructure:"feeModel"`
		Title    string   `mapstructure:"title"`
		Version  int      `mapstructure:"version"`
	} `mapstructure:"AvailableServiceItemList"`
	Ageing string `mapstructure:"Ageing"`
}

type FeeModel struct {
	StartPrice          string `mapstructure:"startPrice" json:"start_price,int"`
	ContinuedHeavy      string `mapstructure:"continuedHeavy" json:"continued_heavy,int"`
	StartWeight         string `mapstructure:"startWeight" json:"start_weight,int"`
	ContinuedHeavyPrice string `mapstructure:"continuedHeavyPrice" json:"continued_heavy_price,int"`
	TotalPrice          string `mapstructure:"totalPrice" json:"total_price,int"`
}
