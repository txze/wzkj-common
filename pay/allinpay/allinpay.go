package allinpay

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hzxiao/goutil"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pay/define"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
)

// PayStrategy 支付策略接口
type PayStrategy interface {
	// Process 处理支付请求
	Process(ctx context.Context, request *common.PaymentRequest) (goutil.Map, error)
	// GetUrl 获取支付接口路径
	GetUrl() string
}

// WxSmallPayStrategy 微信小程序支付策略
type WxSmallPayStrategy struct{}

// NewWxSmallPayStrategy 创建微信小程序支付策略
func NewWxSmallPayStrategy() *WxSmallPayStrategy {
	return &WxSmallPayStrategy{}
}

// Process 处理微信小程序支付请求
func (s *WxSmallPayStrategy) Process(ctx context.Context, request *common.PaymentRequest) (goutil.Map, error) {
	params := goutil.Map{}
	params["goodsName"] = request.GoodsName
	params["outOrderNo"] = request.OrderNo
	params["transAmt"] = strconv.Itoa(request.Amount)
	params["paySource"] = PaySourceDefault
	params["payType"] = PayTypeSmall
	params["funSource"] = "9"
	return params, nil
}

// GetUrl 获取支付接口路径
func (s *WxSmallPayStrategy) GetUrl() string {
	return WxSmallUrl
}

// PayStrategyFactory 支付策略工厂
type PayStrategyFactory struct{}

// NewPayStrategyFactory 创建支付策略工厂
func NewPayStrategyFactory() *PayStrategyFactory {
	return &PayStrategyFactory{}
}

// CreateStrategy 创建支付策略
func (f *PayStrategyFactory) CreateStrategy(productCode string) PayStrategy {
	switch productCode {
	case "WX_SMALL":
		return NewWxSmallPayStrategy()
	default:
		return NewWxSmallPayStrategy()
	}
}

type AllInPay struct {
	config AllInPayConfig
}

func NewAllInPay(cfg AllInPayConfig) (*AllInPay, error) {
	return &AllInPay{
		config: cfg,
	}, nil
}

// Pay 统一支付方法
func (a *AllInPay) Pay(ctx context.Context, request *common.PaymentRequest) (map[string]interface{}, error) {
	// 创建支付策略
	var strategy = NewPayStrategyFactory().CreateStrategy(request.ProductCode)

	// 处理支付请求
	params, err := strategy.Process(ctx, request)
	if err != nil {
		logger.FromContext(ctx).Error("AllInPay Strategy Process Failed", logger.Any("error", err))
		return nil, err
	}

	// 执行请求
	response, err := a.executeRequest(ctx, strategy.GetUrl(), params)
	if err != nil {
		logger.FromContext(ctx).Error("AllInPay Execute Request Failed", logger.Any("error", err))
		return nil, err
	}

	// 根据支付类型返回不同的响应格式
	switch request.ProductCode {
	case "WX_SMALL":
		return response, nil
	default:
		rsp := goutil.Map{}
		rsp["payUrl"] = response["payUrl"]
		return rsp, nil
	}
}

func (a *AllInPay) QueryRefund(ctx context.Context, refundNo, orderNo string) (*common.RefundResponse, error) {
	// 构建查询参数
	params := goutil.Map{}
	params["outOrderNo"] = orderNo
	params["outRefundNo"] = refundNo

	// 执行请求
	response, err := a.executeRequest(ctx, RefundQueryUrl, params)
	if err != nil {
		logger.FromContext(ctx).Error("AllInPay Query Refund Failed", logger.Any("error", err))
		return nil, err
	}

	// 解析响应
	refundAmount, _ := strconv.Atoi(response["refundAmount"].(string))
	successTime := response["successTime"].(string)
	createTime := response["createTime"].(string)

	return &common.RefundResponse{
		UserReceivedAccount:  "",
		SuccessTime:          successTime,
		CreateTime:           createTime,
		RefundStatus:         response.GetString("status") == TradeStatusSuccess,
		OriginalRefundStatus: response.GetString("status"),
		Message:              response.GetString("msg"),
		RefundAmount:         refundAmount,
		Data:                 response,
	}, nil
}

func (a *AllInPay) QueryPayment(ctx context.Context, orderID string) (*common.UnifiedResponse, error) {
	// 构建查询参数
	params := goutil.Map{}
	params["outOrderNo"] = orderID

	// 执行请求
	response, err := a.executeRequest(ctx, QueryUrl, params)
	if err != nil {
		logger.FromContext(ctx).Error("AllInPay Query Payment Failed", logger.Any("error", err))
		return nil, err
	}

	// 解析响应
	amount, _ := strconv.Atoi(response.GetString("transAmt"))
	paidAmount, _ := strconv.Atoi(response.GetString("transAmt"))
	payTime, _ := now.Parse("20060102150405", response.GetString("payTime"))

	return &common.UnifiedResponse{
		Platform:    a.GetType(),
		OrderID:     response.GetString("outOrderNo"),
		PlatformID:  response.GetString("orderNo"),
		Amount:      amount,
		Status:      response.GetString("tradeStatus") == TradeStatusSuccess,
		TradeStatus: a.convertTradeStatus(response.GetString("tradeStatus")),
		PaidAmount:  paidAmount,
		PaidTime:    payTime,
		Message:     response,
	}, nil
}

func (a *AllInPay) Refund(ctx context.Context, request *common.RefundRequest) (*common.RefundOrderResponse, error) {
	// 创建退款请求参数
	refundReq := RefundRequest{
		OutOrderNo:   request.OrderNo,
		OutRefundNo:  request.RefundNo,
		RefundAmount: strconv.Itoa(request.Amount),
	}

	// 构建请求参数
	params := goutil.Map{}
	params["outOrderNo"] = refundReq.OutOrderNo
	params["outRefundNo"] = refundReq.OutRefundNo
	params["refundAmount"] = refundReq.RefundAmount

	// 执行请求
	response, err := a.executeRequest(ctx, a.config.RefundUrl, params)
	if err != nil {
		logger.FromContext(ctx).Error("AllInPay Refund Failed", logger.Any("error", err))
		return nil, err
	}

	// 解析响应
	refundFee, _ := strconv.Atoi(response["refundAmount"].(string))
	successTime := response["successTime"].(string)
	createTime := response["createTime"].(string)

	return &common.RefundOrderResponse{
		OutRefundNo:   request.RefundNo,
		TransactionId: response.GetString("orderNo"),
		OutTradeNo:    request.OrderNo,
		Status:        response.GetString("status"),
		IsSuccess:     response.GetString("status") == TradeStatusSuccess,
		PayerRefund:   refundFee,
		RefundInfo:    response,
		SuccessTime:   successTime,
		CreateTime:    createTime,
	}, nil
}

func (a *AllInPay) Close(ctx context.Context, orderId string) (bool, error) {
	// 构建关闭订单参数
	params := goutil.Map{}
	params["outOrderNo"] = orderId

	// 执行请求
	response, err := a.executeRequest(ctx, RefundUrl, params)
	if err != nil {
		logger.FromContext(ctx).Error("AllInPay Close Order Failed", logger.Any("error", err))
		return false, err
	}

	return response.GetString("status") == TradeStatusSuccess, nil
}

func (a *AllInPay) VerifyNotification(req *http.Request) (*common.UnifiedResponse, error) {
	// 解析请求参数
	bm, err := a.parseNotifyToBodyMap(req)
	if err != nil {
		return nil, err
	}

	// 验证签名
	isCheck, err := a.VerifySign(bm)
	if err != nil {
		return nil, err
	}
	if !isCheck {
		return nil, errors.New("sign is invalid")
	}

	// 解析支付时间
	payTime, _ := now.Parse("20060102150405", bm.GetString("payTime"))

	// 解析金额
	amount, err := strconv.Atoi(bm.GetString("transAmt"))
	if err != nil {
		return nil, ierr.NewIError(ierr.ParamErr, "amount is invalid")
	}

	// 解析交易状态
	tradeStatus := bm.GetString("tradeStatus")

	return &common.UnifiedResponse{
		Platform:    a.GetType(),
		OrderID:     bm.GetString("outOrderNo"),
		PlatformID:  bm.GetString("orderNo"),
		Amount:      amount,
		Status:      tradeStatus == "2", // 2: 支付成功
		TradeStatus: a.convertTradeStatus(tradeStatus),
		PaidAmount:  amount,
		PaidTime:    payTime,
		Message:     bm,
	}, nil
}

func (a *AllInPay) GenerateSign(params map[string]interface{}) (string, error) {
	// 实现签名生成逻辑
	signData := util.Hex(params)
	return util.SM2Sign(a.config.PrivateKey, signData)
}

func (a *AllInPay) VerifySign(params map[string]interface{}) (bool, error) {
	// 实现签名验证逻辑
	sign := params["signature"].(string)
	if sign == "" {
		return false, errors.New("signature is required")
	}
	delete(params, "signature")

	signData := util.Hex(params)
	return util.SM2Verify(a.config.PublicKey, signData, sign)
}

func (a *AllInPay) GetType() string {
	return a.config.GetType()
}

func (a *AllInPay) VerifySettleNotification(ctx context.Context, req *http.Request) (*common.SettleNotificationResponse, error) {
	// 暂时未实现
	return nil, ierr.NewIError(ierr.InternalError, "not implemented")
}

func (a *AllInPay) TradeOrderSettle(ctx context.Context, request common.TradeRoyaltyRateQueryRequestInterface) (*common.TradeRoyaltyRateQueryResponse, error) {
	// 暂时未实现
	return nil, ierr.NewIError(ierr.InternalError, "not implemented")
}

func (a *AllInPay) ConfirmSettle(ctx context.Context, request common.SettleConfirmRequestInterface) (*common.SettleConfirmResponse, error) {
	// 暂时未实现
	return nil, ierr.NewIError(ierr.InternalError, "not implemented")
}

// 执行请求
func (a *AllInPay) executeRequest(ctx context.Context, urlPath string, params goutil.Map) (goutil.Map, error) {
	// 添加基础参数
	params["mchntId"] = a.config.CuSID
	params["storeId"] = a.config.StoreId
	params["channelId"] = a.config.ChannelId
	params["notifyUrl"] = a.config.NotifyUrl
	params["signType"] = a.config.SignType
	params["appId"] = a.config.AppId

	// 生成签名
	signData := util.Hex(params)
	signature, err := util.SM2Sign(a.config.PrivateKey, signData)
	if err != nil {
		return nil, err
	}
	params["signature"] = signature

	// 验证签名
	valid, err := util.SM2Verify(a.config.PublicKey, signData, signature)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid signature")
	}

	// 发送请求
	fullUrl := a.config.PayDomain + urlPath
	response, err := common.Execute(fullUrl, params)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var apiRsp ApiResponse
	if err := json.Unmarshal([]byte(response), &apiRsp); err != nil {
		return nil, errors.Errorf("failed to unmarshal response: %v", err)
	}

	if apiRsp.Code != "0000" {
		return nil, errors.Errorf("api error: %s", apiRsp.Msg)
	}

	// 验证响应签名
	if apiRsp.Sign != "" {
		// 构建验签数据
		signData := util.Hex(apiRsp.Data)
		// 验证签名
		valid, err := util.SM2Verify(a.config.PublicKey, signData, apiRsp.Sign)
		if err != nil {
			return nil, errors.Errorf("failed to verify response signature: %v", err)
		}
		if !valid {
			return nil, errors.New("invalid response signature")
		}
	}

	return apiRsp.Data, nil
}

// 转换交易状态
func (a *AllInPay) convertTradeStatus(status string) string {
	switch status {
	case "2": // 支付成功
		return define.StatusSuccess
	case "1": // 处理中
		return define.StatusPending
	case "3": // 支付失败
		return define.StatusFail
	case "4": // 已关闭
		return define.StatusClose
	default:
		return define.StatusFail
	}
}

// parseNotifyToBodyMap 解析通知请求为 BodyMap
func (a *AllInPay) parseNotifyToBodyMap(req *http.Request) (goutil.Map, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}
	form := req.Form
	bm := goutil.Map{}
	for k, v := range form {
		if len(v) == 1 {
			bm[k] = v[0]
		}
	}
	return bm, nil
}
