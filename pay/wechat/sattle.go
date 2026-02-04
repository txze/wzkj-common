package wechat

import (
	"context"
	"net/http"

	"github.com/go-pay/gopay"
	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/pay/common"
	"github.com/txze/wzkj-common/pkg/util"
)

// 微信分账
type TradeRoyaltyRateQueryRequest struct {
	Appid         string      `json:"appid"`
	SubMchid      string      `json:"sub_mchid"`
	TransactionId string      `json:"transaction_id"`
	OutOrderNo    string      `json:"out_order_no"`
	Receivers     []Receivers `json:"receivers"`
	Finish        bool        `json:"finish"`
}

type Receivers struct {
	Type            string `json:"type"`
	ReceiverAccount string `json:"receiver_account"`
	Amount          int    `json:"amount"`
	Description     string `json:"description"`
	ReceiverName    string `json:"receiver_name"`
}

func (r TradeRoyaltyRateQueryRequest) ToMap() gopay.BodyMap {
	receivers := make([]gopay.BodyMap, 0)
	for _, receiver := range r.Receivers {
		receivers = append(receivers, gopay.BodyMap{
			"type":             receiver.Type,
			"receiver_account": receiver.ReceiverAccount,
			"amount":           receiver.Amount,
			"description":      receiver.Description,
			"receiver_name":    receiver.ReceiverName,
		})
	}
	return gopay.BodyMap{
		"transaction_id": r.TransactionId,
		"out_order_no":   r.OutOrderNo,
		"receivers":      receivers,
		"finish":         r.Finish,
	}
}

// 分账
func (w *Wechat) TradeOrderSettle(ctx context.Context, request common.TradeRoyaltyRateQueryRequestInterface) (*common.TradeRoyaltyRateQueryResponse, error) {
	return nil, nil
}

// VerifySettleNotification 验证分账通知
func (w *Wechat) VerifySettleNotification(ctx context.Context, req *http.Request) (*common.SettleNotificationResponse, error) {
	return nil, nil
}

func (w *Wechat) MapToTradeRoyaltyRateQueryRequest(data goutil.Map) common.TradeRoyaltyRateQueryRequestInterface {
	// 构建请求参数
	receiver := &TradeRoyaltyRateQueryRequest{}
	err := util.S2S(data, receiver)
	if err != nil {
		return nil
	}
	return receiver
}

func (a *Wechat) MapToSettleConfirmRequest(data goutil.Map) common.SettleConfirmRequestInterface {
	return nil
}
