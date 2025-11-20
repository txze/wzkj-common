package logistics

import "github.com/txze/wzkj-common/logistics/model"

// LogisticsProvider 物流接口统一定义
type LogisticsProvider interface {

	//ParseAddress 解析地址信息
	ParseAddress(addr string) (model.Address, error)

	//QueryLogistics 查询物流信息
	QueryLogistics(req *model.QueryLogisticsRequest) (*model.QueryResp, error)

	//CreateOrder 创建物流单
	CreateOrder(req *model.CreateOrderReq) (*model.CreateOrderResp, error)

	//CancelOrder 取消物流单
	CancelOrder(req *model.CancelOrderReq) error

	//ParseWebhook 解析物流推送信息（回调）
	ParseWebhook(body []byte) (*model.WebhookData, error)

	//ParseOrderNotify 订单回调信息
	ParseOrderNotify(body []byte) (*model.OrderNotifyResp, error)

	//GetPriceQuote 获取时效性 以及 预估费用
	GetPriceQuote(req *model.GetPriceQuoteReq) (*model.PriceQuote, error)
}
