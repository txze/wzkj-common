package wechat

import "github.com/go-pay/gopay"

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
