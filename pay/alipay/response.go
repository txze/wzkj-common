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
