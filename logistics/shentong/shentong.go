package shentong

import (
	"fmt"

	"github.com/txze/wzkj-common/logistics/model"
)

type STOClient struct {
	cfg                      *Config
	adaptorCreate            *CreateOrderAdaptor
	adaptorCancel            *CancelOrderAdaptor
	adaptorPriceQuote        *GetPriceQuoteReqAdaptor
	adaptorSubscribeTracking *SubscribeTrackingAdaptor
	adaptorQueryLogistics    *QueryLogisticsAdaptor
}

func (c *STOClient) SubscribeTracking(req *model.SubscribeTrackingReq) error {
	param := c.adaptorSubscribeTracking.ConvertRequest(req)
	formData := convertFormData(STO_TRACE_PLATFORM_SUBSCRIBE, c.cfg.AppKey, c.cfg.ResourceCode, "sto_trace_platform", c.cfg.SecretKey, param)
	baseResp, err := model.DoRequest(c.cfg.GetBaseUrl(), formData)
	if err != nil {
		return err
	}
	return c.adaptorSubscribeTracking.ParseResponse(baseResp)
}

func (c *STOClient) GetPriceQuote(req *model.GetPriceQuoteReq) (*model.PriceQuote, error) {
	param := c.adaptorPriceQuote.ConvertRequest(req)
	formData := convertFormData(QUERY_SEND_SERVICE_DETAIL, c.cfg.AppKey, c.cfg.ResourceCode, "ORDERMS_API", c.cfg.SecretKey, param)
	baseResp, err := model.DoRequest(c.cfg.GetBaseUrl(), formData)
	if err != nil {
		return nil, err
	}
	return c.adaptorPriceQuote.ParseResponse(baseResp)
}

func (c *STOClient) ParseAddress(addr string) (model.Address, error) {
	//TODO implement me
	panic("implement me")
}

func (c *STOClient) QueryLogistics(req *model.QueryLogisticsRequest) (*model.QueryResp, error) {
	param := c.adaptorQueryLogistics.ConvertRequest(req)
	formData := convertFormData(STO_TRACE_QUERY_COMMON, c.cfg.AppKey, c.cfg.ResourceCode, "sto_trace_query", c.cfg.SecretKey, param)
	baseResp, err := model.DoRequest(c.cfg.GetBaseUrl(), formData)
	if err != nil {
		return nil, err
	}

	return c.adaptorQueryLogistics.ParseResponse(req.WaybillNo, baseResp)
}

func (c *STOClient) CreateOrder(req *model.CreateOrderReq) (*model.CreateOrderResp, error) {
	param := c.adaptorCreate.ConvertRequest(c.cfg, req)
	formData := convertFormData(OMS_EXPRESS_ORDER_CREATE, c.cfg.AppKey, c.cfg.ResourceCode, "sto_oms", c.cfg.SecretKey, param)
	baseResp, err := model.DoRequest(c.cfg.GetBaseUrl(), formData)
	if err != nil {
		return nil, err
	}

	return c.adaptorCreate.ParseResponse(baseResp)
}

func (c *STOClient) CancelOrder(req *model.CancelOrderReq) error {
	param := c.adaptorCancel.ConvertRequest(c.cfg, req)
	formData := convertFormData(EDI_MODIFY_ORDER_CANCEL, c.cfg.AppKey, c.cfg.ResourceCode, "edi_modify_order", c.cfg.SecretKey, param)
	baseResp, err := model.DoRequest(c.cfg.GetBaseUrl(), formData)

	if err != nil {
		return err
	}
	return c.adaptorCancel.ParseResponse(baseResp)
}

func (c *STOClient) ParseWebhook(body []byte) (*model.WebhookData, error) {
	fmt.Println("body:", string(body))
	return nil, nil
}

func (c *STOClient) ParseOrderNotify(body []byte) (*model.OrderNotifyResp, error) {
	//TODO implement me
	panic("implement me")
}

func NewSTOClient(cfg *Config) *STOClient {
	return &STOClient{
		cfg:                      cfg,
		adaptorCreate:            &CreateOrderAdaptor{},
		adaptorCancel:            &CancelOrderAdaptor{},
		adaptorPriceQuote:        &GetPriceQuoteReqAdaptor{},
		adaptorSubscribeTracking: &SubscribeTrackingAdaptor{},
		adaptorQueryLogistics:    &QueryLogisticsAdaptor{},
	}
}
