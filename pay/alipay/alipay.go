package alipay

import (
	"context"
	"net/http"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pkg/ierr"
)

type Alipay struct {
	client *alipay.Client
	config AlipayConfig
}

func (a *Alipay) QueryRefund(ctx context.Context, refundNo, orderNo string) (*common.RefundResponse, error) {
	bm := make(gopay.BodyMap)
	bm.
		Set("out_trade_no", orderNo).
		Set("out_request_no", refundNo).
		Set("query_options", []string{
			"deposit_back_info",
			"gmt_refund_pay",
		})
	aliRsp, err := a.client.TradeFastPayRefundQuery(ctx, bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			logger.FromContext(ctx).Error("alipay query refund ", logger.Any("error", bizErr))
			// do something
			return nil, ierr.NewIError(ierr.InternalError, bizErr.Error())
		}
		return nil, err
	}
	logger.FromContext(ctx).Info("alipay query refund ", logger.Any("aliRsp", aliRsp))
	return &common.RefundResponse{
		UserReceivedAccount: aliRsp.Response.DepositBackInfo.EstBankReceiptTime,
		SuccessTime:         aliRsp.Response.DepositBackInfo.BankAckTime,
		CreateTime:          aliRsp.Response.GmtRefundPay,
	}, nil
}

func (a *Alipay) Refund(ctx context.Context, request *common.RefundRequest) error {
	// 请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", request.OrderNo).
		Set("refund_amount", request.Amount).
		Set("refund_reason", request.GoodsName).
		Set("out_request_no", request.RefundNo)

	// 发起退款请求
	aliRsp, err := a.client.TradeRefund(ctx, bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			logger.FromContext(ctx).Error("alipay refund ", logger.Any("error", bizErr))
			// do something
			return err
		}
		return err
	}

	logger.FromContext(ctx).Info("alipay refund success", logger.Any("data", *aliRsp))

	return nil
}

func (a *Alipay) Pay(ctx context.Context, request *common.PaymentRequest) (map[string]interface{}, error) {
	result := decimal.NewFromInt(int64(request.Amount)).Div(decimal.NewFromInt(100))
	//配置公共参数
	a.client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetNotifyUrl(a.config.NotifyUrl)

	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", request.GoodsName)
	bm.Set("out_trade_no", request.OrderNo)
	bm.Set("total_amount", result.String())
	bm.Set("passback_params", request.Params)
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
	totalAmount, err := decimal.NewFromString(bm.GetString("total_amount"))
	if err != nil {
		return nil, err
	}

	buyerPayAmount, err := decimal.NewFromString(bm.GetString("buyer_pay_amount"))
	if err != nil {
		return nil, err
	}

	return &common.UnifiedResponse{
		Platform:    a.GetType(),
		OrderID:     bm.GetString("out_trade_no"),
		PlatformID:  bm.GetString("trade_no"),
		Amount:      int(totalAmount.Mul(decimal.NewFromInt(100)).IntPart()),
		Status:      bm.GetString("trade_status") == "TRADE_SUCCESS",
		TradeStatus: bm.GetString("trade_status"),
		PaidAmount:  int(buyerPayAmount.Mul(decimal.NewFromInt(100)).IntPart()),
		PaidTime:    bm.GetString("gmt_payment"),
		Params:      bm.GetString("passback_params"),
		Message:     bm,
	}, nil
}

func (a *Alipay) QueryPayment(orderID string) (*common.UnifiedResponse, error) {
	//请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderID)

	//查询订单
	aliRsp, err := a.client.TradeQuery(context.Background(), bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			return nil, bizErr
		}
		return nil, err
	}

	totalAmount, err := decimal.NewFromString(bm.GetString("total_amount"))
	if err != nil {
		return nil, err
	}

	buyerPayAmount, err := decimal.NewFromString(bm.GetString("buyer_pay_amount"))
	if err != nil {
		return nil, err
	}

	return &common.UnifiedResponse{
		Platform:    a.GetType(),
		OrderID:     aliRsp.Response.OutTradeNo,
		PlatformID:  aliRsp.Response.TradeNo,
		Amount:      int(totalAmount.Mul(decimal.NewFromInt(100)).IntPart()),
		Status:      aliRsp.Response.TradeStatus == "TRADE_SUCCESS",
		TradeStatus: bm.GetString("trade_status"),
		PaidAmount:  int(buyerPayAmount.Mul(decimal.NewFromInt(100)).IntPart()),
		PaidTime:    aliRsp.Response.SendPayDate,
		Message:     aliRsp,
	}, nil
}

func (a *Alipay) GenerateSign(params map[string]interface{}) (string, error) {
	return "", nil
}

func (a *Alipay) VerifySign(params map[string]interface{}) (bool, error) {
	// 验签
	ok, err := alipay.VerifySignWithCert(a.config.AliPayPublicCertContent, params)
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
	err = client.SetCertSnByContent(cfg.AppCertContent, cfg.AliPayRootCertContent, cfg.AliPayPublicCertContent)
	if err != nil {
		return nil, ierr.NewIError(ierr.InternalError, err.Error())
	}
	return &Alipay{
		client: client,
		config: cfg,
	}, nil
}
