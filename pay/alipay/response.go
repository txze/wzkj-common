package alipay

import "github.com/go-pay/gopay/alipay"

type TradeMergePrecreateResponse struct {
	Response     *TradeMergePrecreate `json:"alipay_trade_merge_precreate_response"`
	Sign         string               `json:"sign"`
	AlipayCertSn string               `json:"alipay_cert_sn,omitempty"`
}

type TradeMergePrecreate struct {
	alipay.ErrorResponse
	OutMergeNo         string                `json:"out_merge_no"`
	PreOrderNo         string                `json:"pre_order_no"`
	OrderDetailResults []*OrderDetailResults `json:"order_detail_results"`
}

type OrderDetailResults struct {
	AppId      string `json:"app_id"`
	OutTradeNo string `json:"out_trade_no"`
	Success    bool   `json:"success"`
	ResultCode string `json:"result_code"`
}

type SettleNotification struct {
	OperationDt         string              `json:"operation_dt"`
	RoyaltyFinishAmount string              `json:"royalty_finish_amount"`
	TradeNo             string              `json:"trade_no"`
	MsgType             string              `json:"msg_type"`
	OperationFinishDt   string              `json:"operation_finish_dt"`
	RoyaltyDetailList   []RoyaltyDetailList `json:"royalty_detail_list"`
	OutRequestNo        string              `json:"out_request_no"`
	SettleNo            string              `json:"settle_no"`
}

type RoyaltyDetailList struct {
	Amount        string `json:"amount"`
	OperationType string `json:"operation_type"`
	TransIn       string `json:"trans_in"`
	State         string `json:"state"`
	DetailId      string `json:"detail_id"`
	ExecuteDt     string `json:"execute_dt"`
	TransInType   string `json:"trans_in_type"`
}
