package wechat

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pay/gopay"
	we "github.com/go-pay/gopay/wechat"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pkg/ierr"
)

type Wechat struct {
	client *wechat.ClientV3
	config WechatConfig
}

func (w *Wechat) Pay(ctx context.Context, request *common.PaymentRequest) (map[string]interface{}, error) {
	//初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("appid", w.config.AppId).
		Set("description", request.GoodsName).
		Set("out_trade_no", request.OrderNo).
		Set("time_expire", request.Expire).
		Set("notify_url", request.NotifyUrl).
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
	rsp := make(map[string]interface{})
	rsp["appId"] = w.config.AppId
	rsp["partnerId"] = w.config.Mchid
	rsp["prepayId"] = wxRsp.Response.PrepayId
	rsp["packageValue"] = w.config.PackageValue
	rsp["nonceStr"] = wxRsp.SignInfo.HeaderNonce
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	rsp["timeStamp"] = timeStamp
	rsp["sign"] = we.GetAppPaySign(w.config.AppId, w.config.Mchid, wxRsp.SignInfo.HeaderNonce, wxRsp.Response.PrepayId, we.SignType_HMAC_SHA256, timeStamp, w.config.ApiV3Key)

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

func (w *Wechat) Refund(orderID string, amount float64) error {
	return nil
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
	mchid := viper.GetString("wechat.pay.mch_id")
	serialNo := viper.GetString("wechat.pay.serialNo")
	apiV3Key := viper.GetString("wechat.pay.apiV3Key")
	privateKey := viper.GetString("wechat.pay.privateKey")
	client, err := wechat.NewClientV3(mchid, serialNo, apiV3Key, privateKey)
	if err != nil {
		return nil, ierr.NewIError(ierr.InternalError, err.Error())
	}

	err = client.AutoVerifySignByPublicKey([]byte(viper.GetString("wechat.pay.PUB_KEY")), viper.GetString("wechat.pay.PUB_KEY_ID"))
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
