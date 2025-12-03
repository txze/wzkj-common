package shentong

import (
	"fmt"
	"net/url"

	"github.com/hzxiao/goutil"
	"github.com/jinzhu/now"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
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
	logger.Info("ParseWebhook request body :", logger.Any("body", string(body)))
	// 解析为 url.Values
	values, err := url.ParseQuery(string(body))
	if err != nil {
		logger.Error("ParseWebhook request body err:", logger.Any("body", string(body)))
		return nil, fmt.Errorf("解析请求体失败: %v", err)
	}

	if values.Get("data_digest") == "" {
		logger.Info("ParseWebhook request body :", logger.Any("body", string(body)))
		return nil, ierr.NewIErrorf(ierr.InternalError, "sign not found :%v", values)
	}

	sign := generateSign(values.Get("content"), c.cfg.SecretKey)
	if sign != values.Get("data_digest") {
		logger.Info("ParseWebhook sign fail :", logger.Any("body", string(body)))
		return nil, ierr.NewIErrorf(ierr.InternalError, "sign fail :%v", values)
	}

	content := values.Get("content")
	rsp := goutil.Map{}
	err = util.Json2S(content, &rsp)
	if err != nil {
		logger.Error("ParseWebhook Json2S content err:", logger.Any("body", string(body)))
		return nil, ierr.NewIErrorf(ierr.InternalError, "ParseOrderNotify json fail :%v", values)
	}

	t, err := now.Parse(rsp.GetStringP("trace/opTime"))
	if err != nil {
		logger.Error("ParseWebhook parse trace/opTime err:", logger.Any("body", string(body)))
	}

	var ret = model.WebhookData{
		OrderId:   "",
		WaybillNo: rsp.GetStringP("waybillNo"),
		ScanType:  rsp.GetStringP("trace/scanType"),
		OpTime:    t,
		Data:      rsp,
	}

	return &ret, nil
}

func (c *STOClient) ParseOrderNotify(body []byte) (*model.OrderNotifyResp, error) {
	logger.Info("ParseOrderNotify request body :", logger.Any("body", string(body)))
	// 解析为 url.Values
	values, err := url.ParseQuery(string(body))
	if err != nil {
		logger.Error("ParseOrderNotify request body err:", logger.Any("body", string(body)))
		return nil, fmt.Errorf("解析请求体失败: %v", err)
	}

	if values.Get("data_digest") == "" {
		logger.Info("ParseOrderNotify request body :", logger.Any("body", string(body)))
		return nil, ierr.NewIErrorf(ierr.InternalError, "sign not found :%v", values)
	}

	sign := generateSign(values.Get("content"), c.cfg.SecretKey)
	if sign != values.Get("data_digest") {
		logger.Info("ParseOrderNotify sign fail :", logger.Any("body", string(body)))
		return nil, ierr.NewIErrorf(ierr.InternalError, "sign fail :%v", values)
	}

	content := values.Get("content")
	rsp := goutil.Map{}
	err = util.Json2S(content, &rsp)
	if err != nil {
		logger.Error("ParseOrderNotify Json2S content err:", logger.Any("body", string(body)))
		return nil, ierr.NewIErrorf(ierr.InternalError, "ParseOrderNotify json fail :%v", values)
	}

	notifyRsp := model.OrderNotifyResp{}
	switch rsp.GetString("event") {
	case EventOrderStatus:
		notifyRsp.OrderId = rsp.GetStringP("changeInfo/OrderId")
		notifyRsp.OriginalStatus = rsp.GetStringP("changeInfo/Status")
		notifyRsp.PickupCode = rsp.GetStringP("changeInfo/PrintCode")
		notifyRsp.WaybillNo = rsp.GetStringP("changeInfo/BillCode")
		notifyRsp.UserMobile = rsp.GetStringP("changeInfo/UserMobile")
		notifyRsp.UserName = rsp.GetStringP("changeInfo/UserName")
		notifyRsp.Status = model.OrderStatusAccept
	case EventOrderCancel:
		notifyRsp.OrderId = rsp.GetStringP("cancelInfo/OrderId")
		notifyRsp.OriginalStatus = OrderStatusCancel.ToString()
		notifyRsp.Reason = rsp.GetStringP("cancelInfo/Reason")
		notifyRsp.WaybillNo = rsp.GetStringP("cancelInfo/BillCode")
		notifyRsp.Status = model.OrderStatusCancel
	case EventOrderUpdateFetchTime:
		notifyRsp.OrderId = rsp.GetStringP("modifyInfo/OrderId")
		notifyRsp.FetchEndTime = rsp.GetStringP("modifyInfo/FetchEndTime")
		notifyRsp.FetchStartTime = rsp.GetStringP("modifyInfo/FetchStartTime")
		notifyRsp.WaybillNo = rsp.GetStringP("modifyInfo/BillCode")
		notifyRsp.Status = model.OrderStatusChangeContract
	case EventOrderRefund:
		notifyRsp.OrderId = rsp.GetStringP("returnInfo/OrderId")
		notifyRsp.Reason = rsp.GetStringP("returnInfo/Reason")
		notifyRsp.OriginalStatus = OrderStatusRefund.ToString()
		notifyRsp.Status = model.OrderStatusRefund
	}

	return &notifyRsp, nil
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
