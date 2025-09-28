package wechat

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"

	"github.com/txze/wzkj-common/pay/common"
)

type Wechat struct {
	client *wechat.ClientV3
	config WechatConfig
}

func (w *Wechat) Pay(request *PaymentRequest) (map[string]interface{}, error) {
	//初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("appid", w.config.AppId).
		Set("description", "测试APP支付商品").
		Set("out_trade_no", request.OrderId).
		Set("time_expire", request.Expire).
		Set("notify_url", request.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", request.Amount).
				Set("currency", request.Currency)
		})

	//请求支付下单，成功后得到结果
	wxRsp, err := w.client.V3TransactionApp(context.Background(), bm)
	if err != nil {
		return nil, err
	}
	if wxRsp.Code != 0 {
		return nil, errors.New(wxRsp.Error)
	}
	rsp := make(map[string]interface{})
	rsp["appId"] = w.config.AppId
	rsp["partnerId"] = w.config.Mchid
	rsp["prepayId"] = wxRsp.Response.PrepayId
	rsp["packageValue"] = w.config.PackageValue
	rsp["nonceStr"] = common.RandomString32Custom()
	rsp["timeStamp"] = time.Now().Unix()
	rsp["sign"] = wxRsp.SignInfo.SignBody

	return rsp, err
}

func (w *Wechat) VerifyNotification(req *http.Request) (*common.UnifiedResponse, error) {
	notifyRsp, err := wechat.V3ParseNotify(req)
	if err != nil {
		return nil, err
	}

	wxPublicKeyMap := w.client.WxPublicKeyMap()
	err = notifyRsp.VerifySignByPKMap(wxPublicKeyMap)
	if err != nil {
		return nil, err
	}

	result, err := notifyRsp.DecryptPayCipherText(w.config.ApiV3Key)
	if err != nil {
		return nil, err
	}

	return &common.UnifiedResponse{
		Platform:   w.GetType(),
		OrderID:    result.OutTradeNo,
		PlatformID: result.TransactionId,
		Amount:     result.Amount.Total,
		Status:     result.TradeState,
		PaidAmount: result.Amount.PayerTotal,
		PaidTime:   result.SuccessTime,
		Message:    result,
	}, nil
}

func (w *Wechat) QueryPayment(orderID string) (*common.UnifiedResponse, error) {
	queryOrder, err := w.client.V3TransactionQueryOrder(context.Background(), wechat.OutTradeNo, orderID)
	if err != nil {
		return nil, err
	}
	return &common.UnifiedResponse{
		Platform:   w.GetType(),
		OrderID:    queryOrder.Response.OutTradeNo,
		PlatformID: queryOrder.Response.TransactionId,
		Amount:     queryOrder.Response.Amount.Total,
		Status:     queryOrder.Response.TradeState,
		PaidAmount: queryOrder.Response.Amount.PayerTotal,
		PaidTime:   queryOrder.Response.SuccessTime,
		Message:    queryOrder,
	}, nil
}

func (w Wechat) Refund(orderID string, amount float64) error {
	//TODO implement me
	panic("implement me")
}

func (w Wechat) GenerateSign(params map[string]interface{}) (string, error) {
	return "", nil
}

func (w Wechat) VerifySign(params map[string]interface{}) (bool, error) {
	return true, nil
}

func (w *Wechat) Close(orderId string) (bool, error) {
	wxRsp, err := w.client.V3TransactionCloseOrder(context.Background(), orderId)
	if err != nil {
		return false, err
	}
	if wxRsp.Code != http.StatusNoContent {
		return false, nil
	}
	return true, nil
}

func (w *Wechat) GetType() string {
	return w.config.GetType() + "_app"
}

func NewWechat(client *wechat.ClientV3, cfg WechatConfig) *Wechat {
	return &Wechat{
		client: client,
		config: cfg,
	}
}
