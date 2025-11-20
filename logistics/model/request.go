package model

import (
	"net/url"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
)

type QueryLogisticsRequest struct {
	ExpressCode string `json:"express_code"` // 快递编码
	WaybillNo   string `json:"waybill_no"`   //物流单号
	Phone       string `json:"phone"`        // 收、寄件人的电话号码
}

// 创建订单信息
type CreateOrderReq struct {
	OrderNo        string   `json:"order_no"` //订单号
	Sender         Sender   `json:"sender"`   //寄件人
	Receiver       Receiver `json:"receiver"` //接收人
	Cargo          Cargo    `json:"cargo"`    //物品信息
	FetchBeginTime string   `json:"fetch_begin_time"`
	FetchEndTime   string   `json:"fetch_end_time"`
}

// 取消寄件
type CancelOrderReq struct {
	OrderId   string `json:"order_id"`   //客户单号 如果是快递100 则传入任务ID
	WaybillNo string `json:"waybill_no"` //物流单号
	Remark    string `json:"remark"`     //备注信息
}

type PickupCodeReq struct {
	OrderId   string `json:"order_id"`   //客户单号 如果是快递100 则传入任务ID
	WaybillNo string `json:"waybill_no"` //物流单号
	Remark    string `json:"remark"`     //备注信息
}

type GetPriceQuoteReq struct {
	SendName    string `json:"SendName"`    //发件人姓名
	SendMobile  string `json:"SendMobile"`  //发件人手机号
	SendProv    string `json:"SendProv"`    //发件人省份
	SendCity    string `json:"SendCity"`    //发件人城市
	SendArea    string `json:"SendArea"`    //发件人区县
	SendAddress string `json:"SendAddress"` //发件人详细地址
	RecName     string `json:"RecName"`     //收件人姓名
	RecMobile   string `json:"RecMobile"`   //收件人手机号
	RecProv     string `json:"RecProv"`     //收件人省份
	RecCity     string `json:"RecCity"`     //收件人城市
	RecArea     string `json:"RecArea"`     //收件人区县
	RecAddress  string `json:"RecAddress"`  //收件人详细地址
	OpenId      string `json:"OpenId"`      //下单用户唯一标识
	Weight      string `json:"Weight"`      //物品重量（kg）
}

type Cargo struct {
	Battery    string `json:"battery"`    //带电标识 （10/未知 20/带电 30/不带电）
	GoodsType  string `json:"goodsType"`  //物品类型（大件、小件、扁平件\文件）
	GoodsName  string `json:"goodsName"`  //物品名称
	GoodsCount int    `json:"goodsCount"` //物品数量
	Weight     int    `json:"weight"`     //kg
}

type Sender struct {
	Name     string `json:"name"`     //寄件人名称
	Tel      string `json:"tel"`      //寄件人固定电话
	Mobile   string `json:"mobile"`   //寄件人手机号码
	PostCode string `json:"postCode"` //邮编
	Country  string `json:"country"`  //国家
	Province string `json:"province"` //省
	City     string `json:"city"`     //市
	Area     string `json:"area"`     //区
	Town     string `json:"town"`     //镇
	Address  string `json:"address"`  //详细地址
}

type Receiver struct {
	Name     string `json:"name"`     //收件人名称
	Tel      string `json:"tel"`      //收件人固定电话
	Mobile   string `json:"mobile"`   //收件人手机号码
	PostCode string `json:"postCode"` //邮编
	Country  string `json:"country"`  //国家
	Province string `json:"province"` //省
	City     string `json:"city"`     //市
	Area     string `json:"area"`     //区
	Town     string `json:"town"`     //镇
	Address  string `json:"address"`  //详细地址
}

// DoRequest 执行API请求
func DoRequest(url string, formData url.Values) (goutil.Map, error) {
	resp, err := util.HttpFormDataPost(url, formData)
	if err != nil {
		return nil, ierr.NewIErrorf(ierr.InternalError, "请求失败: %w", err)
	}

	return resp, nil
}
