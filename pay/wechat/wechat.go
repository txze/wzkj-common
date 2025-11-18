package wechat

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"go.uber.org/zap"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pkg/ierr"
)

type Wechat struct {
	client *wechat.ClientV3
	config WechatConfig
}

func (w *Wechat) QueryRefund(ctx context.Context, refundNo, orderNo string) (*common.RefundResponse, error) {
	bm := make(gopay.BodyMap)
	wxRsp, err := w.client.V3RefundQuery(ctx, refundNo, bm)
	logger.FromContext(ctx).Info("Wechat Query Refund called", zap.Any("wxRsp", wxRsp))
	if err != nil {
		logger.FromContext(ctx).Error("Wechat Query Refund Failed", zap.Error(err))
		return nil, err
	}

	if wxRsp.Code != 0 {
		return nil, ierr.NewIError(ierr.InvalidState, wxRsp.ErrResponse.Message)
	}

	return &common.RefundResponse{
		UserReceivedAccount:  wxRsp.Response.UserReceivedAccount,
		SuccessTime:          wxRsp.Response.SuccessTime,
		CreateTime:           wxRsp.Response.CreateTime,
		OriginalRefundStatus: wxRsp.Response.Status,
		RefundStatus:         wxRsp.Response.Status == gopay.SUCCESS,
		Message:              wxRsp.Error,
	}, nil
}

func (w *Wechat) Pay(ctx context.Context, request *common.PaymentRequest) (map[string]interface{}, error) {
	//初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("appid", w.config.AppId).
		Set("description", request.GoodsName).
		Set("out_trade_no", request.OrderNo).
		Set("time_expire", request.Expire).
		Set("notify_url", w.config.NotifyUrl).
		Set("attach", request.Params).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", request.Amount).
				Set("currency", request.Currency)
		})

	//请求支付下单，成功后得到结果
	wxRsp, err := w.client.V3TransactionApp(context.Background(), bm)
	logger.FromContext(ctx).Info("Wechat Pay called", zap.Any("wxRsp", wxRsp))
	if err != nil {
		logger.FromContext(ctx).Error("Wechat Pay Failed", zap.Error(err))
		return nil, err
	}
	if wxRsp.Code != 0 {
		logger.FromContext(ctx).Error("Wechat Pay Failed", zap.Int("wx_rsp.code", wxRsp.Code))
		return nil, errors.New(wxRsp.Error)
	}
	appPamrams, err := w.client.PaySignOfApp(w.config.AppId, wxRsp.Response.PrepayId)
	if err != nil {
		logger.FromContext(ctx).Error("Wechat Pay Failed", zap.Error(err))
		return nil, err
	}
	rsp := make(map[string]interface{})
	rsp["appId"] = appPamrams.Appid
	rsp["partnerId"] = appPamrams.Partnerid
	rsp["prepayId"] = appPamrams.Prepayid
	rsp["packageValue"] = appPamrams.Package
	rsp["nonceStr"] = appPamrams.Noncestr
	rsp["sign"] = appPamrams.Sign
	rsp["timestamp"] = appPamrams.Timestamp

	return rsp, err
}

func (w *Wechat) VerifyNotification(req *http.Request) (*common.UnifiedResponse, error) {
	notifyRsp, err := wechat.V3ParseNotify(req)
	if err != nil {
		return nil, err
	}

	result, err := notifyRsp.DecryptPayCipherText(string(w.client.ApiV3Key))
	if err != nil {
		return nil, err
	}

	return &common.UnifiedResponse{
		Platform:    w.GetType(),
		OrderID:     result.OutTradeNo,
		PlatformID:  result.TransactionId,
		Amount:      result.Amount.Total,
		Status:      result.TradeState == gopay.SUCCESS,
		TradeStatus: result.TradeState,
		PaidAmount:  result.Amount.PayerTotal,
		PaidTime:    result.SuccessTime,
		Params:      result.Attach,
		Message:     result,
	}, nil
}

func (w *Wechat) QueryPayment(orderID string) (*common.UnifiedResponse, error) {
	queryOrder, err := w.client.V3TransactionQueryOrder(context.Background(), wechat.OutTradeNo, orderID)
	if err != nil {
		return nil, err
	}
	return &common.UnifiedResponse{
		Platform:    w.GetType(),
		OrderID:     queryOrder.Response.OutTradeNo,
		PlatformID:  queryOrder.Response.TransactionId,
		Amount:      queryOrder.Response.Amount.Total,
		Status:      queryOrder.Response.TradeState == gopay.SUCCESS,
		TradeStatus: queryOrder.Response.TradeState,
		PaidAmount:  queryOrder.Response.Amount.PayerTotal,
		PaidTime:    queryOrder.Response.SuccessTime,
		Message:     queryOrder,
	}, nil
}

func (w *Wechat) Refund(ctx context.Context, request *common.RefundRequest) (*common.RefundOrderResponse, error) {
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	// 必填 退款订单号（程序员定义的）
	bm.
		Set("out_refund_no", request.RefundNo).
		Set("out_trade_no", request.OrderNo).
		// 选填 退款描述
		Set("reason", request.GoodsName).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			// 退款金额:单位是分
			bm.Set("refund", request.Amount). //实际退款金额
								Set("total", request.Amount). // 折扣前总金额（不是实际退款数）
								Set("currency", "CNY")
		})
	refund, err := w.client.V3Refund(ctx, bm)
	if err != nil {
		logger.FromContext(ctx).Error("Wechat Refund Failed", zap.Error(err))
		return nil, err
	}

	// 将非正常退款异常记录
	// 返回：404 > {"code":"RESOURCE_NOT_EXISTS","message":"订单不存在"}
	if refund.Code == http.StatusNotFound || refund.Code == http.StatusBadRequest || refund.Code == http.StatusForbidden {
		logger.FromContext(ctx).Error("Wechat Refund Failed", zap.Any("refund", refund))
		return nil, errors.New(refund.Error)
	}
	logger.FromContext(ctx).Info("wechat refund success", logger.Any("data", refund))
	return &common.RefundOrderResponse{
		OutRefundNo:         refund.Response.OutRefundNo,
		TransactionId:       refund.Response.TransactionId,
		OutTradeNo:          refund.Response.OutTradeNo,
		Channel:             refund.Response.Channel,
		UserReceivedAccount: refund.Response.UserReceivedAccount,
		SuccessTime:         refund.Response.SuccessTime,
		CreateTime:          refund.Response.CreateTime,
		Status:              refund.Response.Status,
		Total:               refund.Response.Amount.Total,
		Refund:              refund.Response.Amount.Refund,
		PayerTotal:          refund.Response.Amount.PayerTotal,
		PayerRefund:         refund.Response.Amount.PayerRefund,
		RefundInfo:          refund,
	}, nil
}

func (w *Wechat) GenerateSign(params map[string]interface{}) (string, error) {
	return "", nil
}

func (w *Wechat) VerifySign(params map[string]interface{}) (bool, error) {
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
	return w.config.GetType()
}

func NewWechat(cfg WechatConfig) (*Wechat, error) {
	mchid := cfg.Mchid
	serialNo := cfg.SerialNo
	apiV3Key := cfg.ApiV3Key
	privateKey := cfg.PrivateKey
	client, err := wechat.NewClientV3(mchid, serialNo, apiV3Key, privateKey)
	if err != nil {
		return nil, ierr.NewIError(ierr.InternalError, err.Error())
	}

	err = client.AutoVerifySignByPublicKey([]byte(cfg.PublicKey), cfg.PublicKeyID)
	if err != nil {
		return nil, ierr.NewIError(ierr.InternalError, err.Error())
	}

	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOff

	return &Wechat{
		client: client,
		config: cfg,
	}, nil
}
