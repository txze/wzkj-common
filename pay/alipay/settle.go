package alipay

import (
	"context"
	"net/http"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/hzxiao/goutil"
	"github.com/shopspring/decimal"

	"github.com/jinzhu/copier"

	"github.com/txze/wzkj-common/logger"
	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pay/define"
	"github.com/txze/wzkj-common/pkg/util"
)

// 将金额转换为分
func amountToCents(amountStr string) (int, error) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return 0, err
	}
	return int(amount.Mul(decimal.NewFromInt(100)).IntPart()), nil
}

// SettleDetailInfo 结算详细信息
type SettleDetailInfo struct {
	Amount           int    `json:"amount"`             // 结算的金额，单位为元
	TransIn          string `json:"trans_in"`           // 结算收款方
	TransInType      string `json:"trans_in_type"`      // 结算收款方的账户类型
	SettleEntityID   string `json:"settle_entity_id"`   // 结算主体标识
	SettleEntityType string `json:"settle_entity_type"` // 结算主体类型
	SummaryDimension string `json:"summary_dimension"`  // 结算汇总维度
	ActualAmount     int    `json:"actual_amount"`      // 实际结算金额
}

// SettleInfo 结算信息
type SettleInfo struct {
	SettleDetailInfos []SettleDetailInfo `json:"settle_detail_infos"` // 结算详细信息
}

// SettleConfirmExtendParams 扩展字段信息
type SettleConfirmExtendParams struct {
	RoyaltyFreeze bool `json:"royalty_freeze,omitempty"` // 是否进行资金冻结，用于后续分账
	RoyaltyFinish bool `json:"royalty_finish,omitempty"` // 是否完结分账，默认false
}

// settleConfirmRequest 结算确认请求
type settleConfirmRequest struct {
	OutRequestNo string                     `json:"out_request_no"` // 确认结算请求流水号
	TradeNo      string                     `json:"trade_no"`       // 支付宝交易号
	SettleInfo   *SettleInfo                `json:"settle_info"`    // 结算信息
	ExtendParams *SettleConfirmExtendParams `json:"extend_params"`  // 扩展字段信息
}

// GetPlatform 实现SettleConfirmRequestInterface接口
func (r *settleConfirmRequest) GetPlatform() string {
	return define.PlatformAlipay
}

// ToStruct 将任意类型的结构体转换为指定类型的结构体
func (r *settleConfirmRequest) ToStruct(to interface{}) error {
	err := copier.Copy(to, r)
	if err != nil {
		return err
	}
	return nil
}

// SettleConfirmResponse 结算确认响应
type SettleConfirmResponse struct {
	Code    string `json:"code"`     // 响应码
	Msg     string `json:"msg"`      // 响应信息
	SubCode string `json:"sub_code"` // 子响应码
	SubMsg  string `json:"sub_msg"`  // 子响应信息
}

// Confirm 资金结算确认请求
func (a *Alipay) confirm(ctx context.Context, request *settleConfirmRequest) (*SettleConfirmResponse, error) {
	// 配置公共参数
	a.client.SetCharset("utf-8").
		SetSignType(alipay.RSA2)

	// 构建请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_request_no", request.OutRequestNo)
	bm.Set("trade_no", request.TradeNo)

	// 构建结算信息
	if request.SettleInfo != nil {
		settleDetailInfos := make([]gopay.BodyMap, len(request.SettleInfo.SettleDetailInfos))
		for i, detail := range request.SettleInfo.SettleDetailInfos {
			settleDetailInfo := make(gopay.BodyMap)
			if detail.Amount > 0 {
				settleDetailInfo.Set("amount", centsToAmount(int64(detail.Amount)))
			}
			settleDetailInfo.Set("trans_in", detail.TransIn)
			settleDetailInfo.Set("trans_in_type", detail.TransInType)
			if detail.SettleEntityID != "" {
				settleDetailInfo.Set("settle_entity_id", detail.SettleEntityID)
			}
			if detail.SettleEntityType != "" {
				settleDetailInfo.Set("settle_entity_type", detail.SettleEntityType)
			}
			if detail.SummaryDimension != "" {
				settleDetailInfo.Set("summary_dimension", detail.SummaryDimension)
			}
			if detail.ActualAmount > 0 {
				settleDetailInfo.Set("actual_amount", centsToAmount(int64(detail.ActualAmount)))
			}
			settleDetailInfos[i] = settleDetailInfo
		}

		settleInfo := make(gopay.BodyMap)
		settleInfo.Set("settle_detail_infos", settleDetailInfos)
		bm.Set("settle_info", settleInfo)
	}

	// 添加扩展参数
	if request.ExtendParams != nil {
		extendParams := make(gopay.BodyMap)
		extendParams.Set("royalty_freeze", request.ExtendParams.RoyaltyFreeze)
		bm.Set("extend_params", extendParams)
	}

	// 发起结算确认请求
	aliRsp, err := a.client.TradeSettleConfirm(ctx, bm)
	if err != nil {
		logger.FromContext(ctx).Error("alipay settle confirm error", logger.Any("error", err))
		return nil, a.handleBizError(ctx, err, "alipay trade settle confirm")
	}

	logger.FromContext(ctx).Info("alipay settle confirm success", logger.Any("response", aliRsp))

	// 构建响应
	response := &SettleConfirmResponse{
		Code:    aliRsp.Response.Code,
		Msg:     aliRsp.Response.Msg,
		SubCode: aliRsp.Response.SubCode,
		SubMsg:  aliRsp.Response.SubMsg,
	}

	return response, nil
}

// 支付宝分账
type TradeRoyaltyRateQueryRequest struct {
	OutRequestNo      string `json:"out_request_no"`
	TradeNo           string `json:"trade_no"`
	RoyaltyParameters []struct {
		RoyaltyType      string `json:"royalty_type"`
		TransInType      string `json:"trans_in_type"`
		TransIn          string `json:"trans_in"`
		Amount           int    `json:"amount"`
		AmountPercentage int    `json:"amount_percentage"`
		Desc             string `json:"desc"`
	} `json:"royalty_parameters"`
	RoyaltyMode  string                     `json:"royalty_mode"`  // 分账模式
	ExtendParams *SettleConfirmExtendParams `json:"extend_params"` // 扩展字段信息
}

func (a *Alipay) MapToTradeRoyaltyRateQueryRequest(data goutil.Map) common.TradeRoyaltyRateQueryRequestInterface {
	// 构建请求参数
	receiver := &TradeRoyaltyRateQueryRequest{}
	err := util.S2S(data, receiver)
	if err != nil {
		return nil
	}
	return receiver
}

func (a *Alipay) MapToSettleConfirmRequest(data goutil.Map) common.SettleConfirmRequestInterface {
	// 构建请求参数
	receiver := &settleConfirmRequest{}
	err := util.S2S(data, receiver)
	if err != nil {
		return nil
	}
	return receiver
}

func (receiver TradeRoyaltyRateQueryRequest) ToMap() gopay.BodyMap {
	bm := make(gopay.BodyMap)
	bm.Set("out_request_no", receiver.OutRequestNo)
	bm.Set("trade_no", receiver.TradeNo)
	royaltyParameters := make([]gopay.BodyMap, 0)
	for _, parameter := range receiver.RoyaltyParameters {
		royaltyParameter := make(gopay.BodyMap)
		royaltyParameter.Set("royalty_type", parameter.RoyaltyType)
		royaltyParameter.Set("trans_in_type", parameter.TransInType)
		royaltyParameter.Set("trans_in", parameter.TransIn)
		if parameter.AmountPercentage > 0 {
			royaltyParameter.Set("amount_percentage", parameter.AmountPercentage)
		} else {
			royaltyParameter.Set("amount", centsToAmount(int64(parameter.Amount)))
		}
		royaltyParameter.Set("desc", parameter.Desc)
		royaltyParameters = append(royaltyParameters, royaltyParameter)
	}
	bm.Set("royalty_parameters", royaltyParameters)
	if receiver.RoyaltyMode == "" {
		receiver.RoyaltyMode = define.RoyaltyModeAsync
	}
	bm.Set("royalty_mode", receiver.RoyaltyMode)

	if receiver.ExtendParams != nil {
		bm.Set("extend_params", gopay.BodyMap{
			"royalty_finish": receiver.ExtendParams.RoyaltyFinish,
		})
	}

	return bm
}

// 分账
func (a *Alipay) TradeOrderSettle(ctx context.Context, request common.TradeRoyaltyRateQueryRequestInterface) (*common.TradeRoyaltyRateQueryResponse, error) {
	// 配置公共参数
	a.client.SetCharset("utf-8").
		SetSignType(alipay.RSA2)

	// 构建请求参数
	bm := request.ToMap()

	// 发起分账比例查询请求
	aliRsp, err := a.client.TradeOrderSettle(ctx, bm)
	if err != nil {
		logger.FromContext(ctx).Error("alipay royalty rate query error", logger.Any("error", err))
		return nil, a.handleBizError(ctx, err, "alipay trade royalty rate query")
	}

	logger.FromContext(ctx).Info("alipay royalty rate query success", logger.Any("response", aliRsp))

	// 构建响应alipay_trade_order_settle_response
	response := &common.TradeRoyaltyRateQueryResponse{
		Platform:            define.PlatformAlipay,
		TransactionId:       aliRsp.Response.TradeNo,
		OutOrderNo:          bm.GetString("out_request_no"),
		TransactionSettleNo: aliRsp.Response.SettleNo,
		Code:                aliRsp.Response.Code,
		Msg:                 aliRsp.Response.Msg,
		RawData:             aliRsp,
	}

	return response, nil
}

// VerifySettleNotification 验证分账通知
func (a *Alipay) VerifySettleNotification(ctx context.Context, req *http.Request) (*common.SettleNotificationResponse, error) {
	// 验证分账通知
	// 解析请求参数
	bm, err := alipay.ParseNotifyToBodyMap(req)
	if err != nil {
		return nil, err
	}

	_, err = a.VerifySign(bm)
	if err != nil {
		return nil, err
	}

	// 解析分账通知参数
	outRequestNo := bm.GetString("out_request_no")
	msgType := bm.GetString("msg_type")
	tradeNo := bm.GetString("trade_no")
	royaltyFinishAmountStr := bm.GetString("royalty_finish_amount")
	settleNo := bm.GetString("settle_no")
	operationDt := bm.GetString("operation_dt")
	operationFinishDt := bm.GetString("operation_finish_dt")

	// 转换分账完结金额为分
	var royaltyFinishAmount int
	if royaltyFinishAmountStr != "" {
		royaltyFinishAmount, err = amountToCents(royaltyFinishAmountStr)
		if err != nil {
			logger.FromContext(ctx).Error("alipay verify settle notification error", logger.Any("error", err))
			return nil, err
		}
	}

	// 解析分账明细
	var royaltyDetailList []common.RoyaltyDetail
	royaltyDetailListStr := bm.GetString("royalty_detail_list")
	if royaltyDetailListStr != "" {
		var royaltyDetails []struct {
			OperationType  string `json:"operation_type"`
			Amount         string `json:"amount"`
			State          string `json:"state"`
			ExecuteDt      string `json:"execute_dt"`
			TransOut       string `json:"trans_out"`
			TransOutType   string `json:"trans_out_type"`
			TransOutOpenId string `json:"trans_out_open_id"`
			TransIn        string `json:"trans_in"`
			TransInType    string `json:"trans_in_type"`
			TransInOpenId  string `json:"trans_in_open_id"`
			DetailId       string `json:"detail_id"`
			ErrorCode      string `json:"error_code"`
			ErrorDesc      string `json:"error_desc"`
		}

		err = util.Json2S(royaltyDetailListStr, &royaltyDetails)
		if err != nil {
			logger.FromContext(ctx).Error("alipay verify settle notification error", logger.Any("error", err))
			return nil, err
		}

		for _, detail := range royaltyDetails {
			amount, err := amountToCents(detail.Amount)
			if err != nil {
				logger.FromContext(ctx).Error("alipay verify settle notification error", logger.Any("error", err))
				return nil, err
			}

			royaltyDetail := common.RoyaltyDetail{
				OperationType:  detail.OperationType,
				Amount:         amount,
				State:          detail.State,
				ExecuteDt:      detail.ExecuteDt,
				TransOut:       detail.TransOut,
				TransOutType:   detail.TransOutType,
				TransOutOpenId: detail.TransOutOpenId,
				TransIn:        detail.TransIn,
				TransInType:    detail.TransInType,
				TransInOpenId:  detail.TransInOpenId,
				DetailId:       detail.DetailId,
				ErrorCode:      detail.ErrorCode,
				ErrorDesc:      detail.ErrorDesc,
			}
			royaltyDetailList = append(royaltyDetailList, royaltyDetail)
		}
	}

	// 构建分账通知响应
	response := &common.SettleNotificationResponse{
		Platform:            a.GetType(),
		OutRequestNo:        outRequestNo,
		MsgType:             msgType,
		TradeNo:             tradeNo,
		RoyaltyFinishAmount: royaltyFinishAmount,
		SettleNo:            settleNo,
		OperationDt:         operationDt,
		OperationFinishDt:   operationFinishDt,
		RoyaltyDetailList:   royaltyDetailList,
		RawData:             bm,
	}

	logger.FromContext(ctx).Info("alipay verify settle notification success", logger.Any("response", response))
	return response, nil
}
