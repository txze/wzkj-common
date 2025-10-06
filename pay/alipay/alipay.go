package alipay

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/pkg/errors"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pkg/ierr"
)

type Alipay struct {
	client *alipay.Client
	config AlipayConfig
}

func (a *Alipay) Refund(ctx context.Context, request *common.RefundRequest) error {
	return nil
}

func (a *Alipay) Pay(ctx context.Context, request *common.PaymentRequest) (map[string]interface{}, error) {
	//配置公共参数
	a.client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetNotifyUrl(a.config.NotifyUrl)

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", request.GoodsName)
	bm.Set("out_trade_no", request.OrderNo)
	bm.Set("total_amount", request.Amount)
	//手机APP支付参数请求
	payParam, err := a.client.TradeAppPay(context.Background(), bm)
	if err != nil {
		logger.FromContext(ctx).Error("alipay error: " + err.Error())
		return nil, err
	}
	rsp := make(map[string]interface{})
	rsp["orderStr"] = payParam
	return rsp, err
}

func (a *Alipay) VerifyNotification(req *http.Request) (*common.UnifiedResponse, error) {
	// 解析请求参数
	bm, err := alipay.ParseNotifyToBodyMap(req)
	if err != nil {
		return nil, err
	}
	_, err = a.VerifySign(bm)
	if err != nil {
		return nil, err
	}
	totalAmount, _ := strconv.ParseFloat(bm.GetString("total_amount"), 64)
	buyerPayAmount, _ := strconv.ParseFloat(bm.GetString("buyer_pay_amount"), 64)
	return &common.UnifiedResponse{
		Platform:   a.GetType(),
		OrderID:    bm.GetString("out_trade_no"),
		PlatformID: bm.GetString("trade_no"),
		Amount:     totalAmount,
		Status:     bm.GetString("trade_status"),
		PaidAmount: buyerPayAmount,
		PaidTime:   bm.GetString("gmt_payment"),
		Message:    bm,
	}, nil
}

func (a *Alipay) QueryPayment(orderID string) (*common.UnifiedResponse, error) {
	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderID)

	//查询订单
	aliRsp, err := a.client.TradeQuery(context.Background(), bm)
	if err != nil {
		return nil, err
	}
	if aliRsp.Response.Code != "10000" {
		return nil, errors.New(aliRsp.Response.Msg)
	}
	totalAmount, _ := strconv.ParseFloat(aliRsp.Response.TotalAmount, 64)
	buyerPayAmount, _ := strconv.ParseFloat(aliRsp.Response.BuyerPayAmount, 64)
	return &common.UnifiedResponse{
		Platform:   a.GetType(),
		OrderID:    aliRsp.Response.OutTradeNo,
		PlatformID: aliRsp.Response.TradeNo,
		Amount:     totalAmount,
		Status:     aliRsp.Response.TradeStatus,
		PaidAmount: buyerPayAmount,
		PaidTime:   aliRsp.Response.SendPayDate,
		Message:    aliRsp,
	}, nil
}

func (a *Alipay) GenerateSign(params map[string]interface{}) (string, error) {
	return "", nil
}

func (a *Alipay) VerifySign(params map[string]interface{}) (bool, error) {
	// 验签
	ok, err := alipay.VerifySign(a.config.AliPayPublicKey, params)
	if err != nil {
		return false, errors.Errorf("alipay.VerifySign err: %s", err.Error())
	}
	if ok == false {
		return false, errors.New("alipay.VerifySign err")
	}
	return ok, nil
}

func (a *Alipay) Close(orderId string) (bool, error) {
	// 请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderId)

	// 撤销支付订单
	aliRsp, err := a.client.TradeClose(context.Background(), bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			return false, bizErr
		}
		return false, err
	}
	if aliRsp.Response.Code != "10000" {
		return false, errors.Errorf("alipay error: %s", aliRsp.Response.Msg)
	}
	return true, nil
}

func (a *Alipay) GetType() string {
	return a.config.GetType()
}

func NewAlipay(cfg AlipayConfig) (*Alipay, error) {
	client, err := alipay.NewClient(cfg.Appid, cfg.PrivateKey, cfg.IsProd)
	if err != nil {
		return nil, ierr.NewIError(ierr.InternalError, err.Error())
	}
	return &Alipay{
		client: client,
		config: cfg,
	}, nil
}
