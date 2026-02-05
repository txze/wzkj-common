package alipay

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/hzxiao/goutil"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pay/define"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
)

type Alipay struct {
	client *alipay.Client
	config AlipayConfig
}

// 处理支付宝业务错误
func (a *Alipay) handleBizError(ctx context.Context, err error, operation string) error {
	if bizErr, ok := alipay.IsBizError(err); ok {
		logger.FromContext(ctx).Error(operation, logger.Any("error", bizErr))
		return ierr.NewIError(ierr.InternalError, bizErr.Error())
	}
	return err
}

// 将分转换为元
func centsToAmount(cents int64) string {
	return decimal.NewFromInt(cents).Div(decimal.NewFromInt(100)).String()
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
		return nil, a.handleBizError(ctx, err, "alipay query refund")
	}

	logger.FromContext(ctx).Info("alipay query refund", logger.Any("aliRsp", aliRsp))

	successTime := aliRsp.Response.GmtRefundPay
	createTime := ""
	if aliRsp.Response.DepositBackInfo.EstBankReceiptTime != "" {
		successTime = aliRsp.Response.DepositBackInfo.EstBankReceiptTime
		createTime = aliRsp.Response.GmtRefundPay
	}

	refundAmountInt, err := amountToCents(aliRsp.Response.RefundAmount)
	if err != nil {
		return nil, err
	}

	return &common.RefundResponse{
		UserReceivedAccount:  "",
		SuccessTime:          successTime,
		CreateTime:           createTime,
		RefundStatus:         aliRsp.Response.RefundStatus == "REFUND_SUCCESS",
		OriginalRefundStatus: aliRsp.Response.RefundStatus,
		Message:              aliRsp.Response.Msg,
		RefundAmount:         refundAmountInt,
		Data:                 aliRsp,
	}, nil
}

func (a *Alipay) Refund(ctx context.Context, request *common.RefundRequest) (*common.RefundOrderResponse, error) {
	// 请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", request.OrderNo).
		Set("refund_amount", centsToAmount(int64(request.Amount))).
		Set("refund_reason", request.GoodsName).
		Set("out_request_no", request.RefundNo)

	// 发起退款请求
	aliRsp, err := a.client.TradeRefund(ctx, bm)
	if err != nil {
		return nil, a.handleBizError(ctx, err, "alipay refund")
	}

	queryRefund, err := a.QueryRefund(ctx, request.RefundNo, request.OrderNo)
	if err != nil {
		logger.FromContext(ctx).Error("alipay query refund", logger.Any("error", err))
		return nil, ierr.NewIError(ierr.InternalError, err.Error())
	}

	if !queryRefund.RefundStatus {
		logger.FromContext(ctx).Error("alipay refund", logger.Any("status", queryRefund.RefundStatus))
		return nil, ierr.NewIError(ierr.InternalError, "退款失败:"+queryRefund.Message)
	}

	logger.FromContext(ctx).Info("alipay refund success", logger.Any("data", *aliRsp))
	refundFeeInt, err := amountToCents(aliRsp.Response.RefundFee)
	if err != nil {
		return nil, err
	}

	return &common.RefundOrderResponse{
		OutRefundNo:         request.RefundNo,
		TransactionId:       aliRsp.Response.TradeNo,
		OutTradeNo:          aliRsp.Response.OutTradeNo,
		Channel:             "",
		UserReceivedAccount: queryRefund.UserReceivedAccount,
		SuccessTime:         queryRefund.SuccessTime,
		CreateTime:          queryRefund.CreateTime,
		Status:              queryRefund.OriginalRefundStatus,
		IsSuccess:           queryRefund.RefundStatus,
		PayerRefund:         refundFeeInt,
		RefundInfo:          aliRsp,
	}, nil
}

// 普通支付/服务商模式 单笔订单支付
func (a *Alipay) Pay(ctx context.Context, request *common.PaymentRequest) (map[string]interface{}, error) {
	// 配置公共参数
	a.client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetNotifyUrl(a.config.NotifyUrl)

	// 请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", request.GoodsName)
	bm.Set("out_trade_no", request.OrderNo)
	bm.Set("total_amount", centsToAmount(int64(request.Amount)))
	bm.Set("passback_params", request.Params)
	if request.IsServiceMode {
		bm.Set("product_code", request.ProductCode)
		bm.Set("sub_merchant", goutil.Map{
			"merchant_id": request.SubMerchantID,
		})
		for i, _ := range request.SettleDetailInfos {
			request.SettleDetailInfos[i]["amount"] = centsToAmount(int64(request.Amount))
		}
		bm.Set("settle_info", goutil.Map{
			"settle_detail_infos": request.SettleDetailInfos,
		})

		bm.Set("settle_period_time", "30d")
	}

	// 手机APP支付参数请求
	payParam, err := a.client.TradeAppPay(ctx, bm)
	if err != nil {
		logger.FromContext(ctx).Error("alipay error", logger.Any("error", err))
		return nil, err
	}

	rsp := make(map[string]interface{})
	rsp["orderStr"] = payParam
	return rsp, nil
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

	logger.FromContext(req.Context()).Info("alipay verify data", logger.Any("data", bm))

	totalAmountInt, err := amountToCents(bm.GetString("total_amount"))
	if err != nil {
		return nil, err
	}

	buyerPayAmountInt, err := amountToCents(bm.GetString("buyer_pay_amount"))
	if err != nil {
		return nil, err
	}

	t, _ := now.Parse(time.DateTime, bm.GetString("gmt_payment"))

	discountAmount := 0
	voucherDetailList := bm.GetString("voucher_detail_list")
	if voucherDetailList != "" {
		var voucherDetailListResp []alipay.NotifyVoucherDetail
		err = util.Json2S(voucherDetailList, &voucherDetailListResp)
		if err != nil {
			logger.FromContext(req.Context()).Error("alipay verify notification error", logger.Any("error", err))
		} else {
			for _, detail := range voucherDetailListResp {
				discountInt, err := amountToCents(detail.Amount)
				if err != nil {
					logger.FromContext(req.Context()).Error("alipay verify notification error", logger.Any("error", err))
				} else {
					discountAmount += discountInt
				}
			}
		}
	}

	return &common.UnifiedResponse{
		Platform:       a.GetType(),
		OrderID:        bm.GetString("out_trade_no"),
		PlatformID:     bm.GetString("trade_no"),
		Amount:         totalAmountInt,
		Status:         bm.GetString("trade_status") == "TRADE_SUCCESS",
		TradeStatus:    common.ConvertPaymentStatus(a.GetType(), bm.GetString("trade_status")),
		PaidAmount:     buyerPayAmountInt,
		PaidTime:       t,
		Params:         bm.GetString("passback_params"),
		Message:        bm,
		DiscountAmount: discountAmount,
	}, nil
}

func (a *Alipay) QueryPayment(ctx context.Context, orderID string) (*common.UnifiedResponse, error) {
	// 请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderID)

	// 查询订单
	aliRsp, err := a.client.TradeQuery(ctx, bm)
	if err != nil {
		return nil, a.handleBizError(ctx, err, "alipay query payment")
	}

	totalAmountInt, err := amountToCents(aliRsp.Response.TotalAmount)
	if err != nil {
		return nil, err
	}

	buyerPayAmountInt, err := amountToCents(aliRsp.Response.BuyerPayAmount)
	if err != nil {
		return nil, err
	}

	t, _ := now.Parse(time.DateTime, aliRsp.Response.SendPayDate)

	return &common.UnifiedResponse{
		Platform:    a.GetType(),
		OrderID:     aliRsp.Response.OutTradeNo,
		PlatformID:  aliRsp.Response.TradeNo,
		Amount:      totalAmountInt,
		Status:      aliRsp.Response.TradeStatus == "TRADE_SUCCESS",
		TradeStatus: aliRsp.Response.TradeStatus,
		PaidAmount:  buyerPayAmountInt,
		PaidTime:    t,
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
	if !ok {
		return false, errors.New("alipay.VerifySign err")
	}
	return ok, nil
}

func (a *Alipay) Close(ctx context.Context, orderId string) (bool, error) {
	// 请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", orderId)

	// 关闭支付订单
	aliRsp, err := a.client.TradeClose(ctx, bm)
	if err != nil {
		return false, a.handleBizError(ctx, err, "alipay close order")
	}

	if aliRsp.Response.Code != "10000" {
		return false, errors.Errorf("alipay error: %s", aliRsp.Response.Msg)
	}

	return true, nil
}

func (a *Alipay) MergePay(ctx context.Context, bm gopay.BodyMap) (goutil.Map, error) {
	// 参数验证
	if bm.GetString("out_merge_no") == gopay.NULL {
		return nil, errors.New("out_merge_no is required")
	}

	// 配置公共参数
	a.client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetNotifyUrl(a.config.NotifyUrl)

	// 发起合单支付请求
	bs, err := a.client.DoAliPay(ctx, bm, "alipay.trade.merge.precreate")
	if err != nil {
		return nil, err
	}

	// 解析响应参数
	var aliRsp *TradeMergePrecreateResponse
	aliRsp = new(TradeMergePrecreateResponse)
	if err = json.Unmarshal(bs, &aliRsp); err != nil || aliRsp.Response == nil {
		return nil, errors.Errorf("[%v], bytes: %s", gopay.UnmarshalErr, string(bs))
	}

	// 检查业务错误
	if aliRsp.Response.Code != "10000" {
		return nil, errors.Errorf("alipay error: %s", aliRsp.Response.Msg)
	}

	// 返回响应数据
	responseMap := goutil.Map{
		"out_merge_no":         aliRsp.Response.OutMergeNo,
		"pre_order_no":         aliRsp.Response.PreOrderNo,
		"order_detail_results": aliRsp.Response.OrderDetailResults,
	}

	return responseMap, nil
}

func (a *Alipay) GetType() string {
	return a.config.GetType()
}

// ConfirmSettle 结算确认
func (a *Alipay) ConfirmSettle(ctx context.Context, request common.SettleConfirmRequestInterface) (*common.SettleConfirmResponse, error) {
	// 将通用接口转换为支付宝特定的请求结构

	var alipayRequest settleConfirmRequest
	err := request.ToStruct(&alipayRequest)
	if err != nil {
		return nil, errors.New("invalid request type for alipay settle confirm")
	}

	// 调用支付宝的Confirm方法
	response, err := a.confirm(ctx, &alipayRequest)
	if err != nil {
		return nil, err
	}

	// 转换为通用响应结构
	return &common.SettleConfirmResponse{
		Platform: define.PlatformAlipay,
		Code:     response.Code,
		Msg:      response.Msg,
		SubCode:  response.SubCode,
		SubMsg:   response.SubMsg,
		RawData:  response,
	}, nil
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
