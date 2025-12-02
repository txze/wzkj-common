package model

import (
	"encoding/json"
	"errors"
)

type QueryResp struct {
	WaybillNo   string  `json:"waybill_no"` //物流单号
	Ischeck     string  `json:"ischeck"`
	ExpressCode string  `json:"express_code"` // 快递编码
	Status      string  `json:"status"`
	ItemStatus  string  `json:"item_status"`
	State       string  `json:"state"`
	Data        []*Data `json:"data"`
}

type Data struct {
	Time       string `json:"time"`
	Context    string `json:"context"`
	Ftime      string `json:"ftime"`
	AreaCode   string `json:"areaCode"`
	AreaName   string `json:"areaName"`
	Status     string `json:"status"`
	Location   string `json:"location"`
	AreaCenter string `json:"areaCenter"`
	AreaPinYin string `json:"areaPinYin"`
	StatusCode string `json:"statusCode"`
}

type CreateOrderResp struct {
	OrderId   string `json:"order_id"`   //订单号（客户系统自己生成，唯一）
	WaybillNo string `json:"waybill_no"` //物流单号
}

func parseJSON[T any](body []byte, target *T) error {
	if len(body) == 0 {
		return errors.New("body is empty")
	}
	return json.Unmarshal(body, target)
}

const (
	OrderStatusAccept         = "已接单"
	OrderStatusCancel         = "已取消" //已取消
	OrderStatusChangeContract = "改约"  //改约
	OrderStatusRefund         = "已下单"
)

// 解析回调订单信息
type OrderNotifyResp struct {
	OrderId        string `json:"order_id"`         //订单ID
	WaybillNo      string `json:"waybill_no"`       //物流单号
	Status         string `json:"status"`           //订单调度状态
	OriginalStatus string `json:"original_status"`  //订单原始调度状态
	UserCode       string `json:"user_code"`        //接单业务员编号
	UserName       string `json:"user_name"`        //接单业务员名称
	UserMobile     string `json:"user_mobile"`      //业务员手机号
	FetchStartTime string `json:"fetch_start_time"` //预约取件开始时间
	FetchEndTime   string `json:"fetch_end_time"`   //预约取件结束时间
	PickupCode     string `json:"pickup_code"`      //取件码
	Reason         string `json:"reason"`           //原因
}

// 解析物流回调信息
type WebhookData struct {
	OrderId   string `json:"order_id"`   //订单ID
	WaybillNo string `json:"waybill_no"` //物流单号
	ScanType  string `json:"scan_type"`  //物流状态
}

type PriceQuote struct {
	StartPrice          int `json:"start_price"`
	ContinuedHeavy      int `json:"continued_heavy"`
	StartWeight         int `json:"start_weight,int"`
	ContinuedHeavyPrice int `json:"continued_heavy_price"`
	TotalPrice          int `json:"total_price"`
}
