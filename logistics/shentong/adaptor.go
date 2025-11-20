package shentong

import (
	"fmt"
	"strconv"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/pkg/ierr"
)

type CreateOrderAdaptor struct {
}

func (c *CreateOrderAdaptor) ConvertRequest(cfg *Config, req *model.CreateOrderReq) *CreateOrderRequest {
	return &CreateOrderRequest{
		OrderNo:     req.OrderNo,
		OrderSource: cfg.SourceCode,
		BillType:    "00",
		OrderType:   "02",
		Sender: Sender{
			Name:     req.Sender.Name,
			Tel:      req.Sender.Tel,
			Mobile:   req.Sender.Mobile,
			PostCode: req.Sender.PostCode,
			Country:  req.Sender.Country,
			Province: req.Sender.Province,
			City:     req.Sender.City,
			Area:     req.Sender.Area,
			Town:     req.Sender.Town,
			Address:  req.Sender.Address,
		},
		Receiver: Receiver{
			Name:     req.Receiver.Name,
			Tel:      req.Receiver.Tel,
			Mobile:   req.Receiver.Mobile,
			PostCode: req.Receiver.PostCode,
			Country:  req.Receiver.Country,
			Province: req.Receiver.Province,
			City:     req.Receiver.City,
			Area:     req.Receiver.Area,
			Town:     req.Receiver.Town,
			Address:  req.Receiver.Address,
		},
		Cargo: Cargo{
			Battery:    req.Cargo.Battery,
			GoodsType:  req.Cargo.GoodsType,
			GoodsName:  req.Cargo.GoodsName,
			GoodsCount: req.Cargo.GoodsCount,
			Weight:     req.Cargo.Weight,
		},
		Customer: Customer{
			SiteCode:          cfg.Customer.SiteCode,
			CustomerName:      cfg.Customer.CustomerName,
			SitePwd:           cfg.Customer.SitePwd,
			MonthCustomerCode: cfg.Customer.MonthCustomerCode,
		},
		ExtendFieldMap: goutil.Map{
			"fetch_begin_time": req.FetchBeginTime,
			"fetch_end_time":   req.FetchEndTime,
		},
	}
}

func (c *CreateOrderAdaptor) ParseResponse(rspMap goutil.Map) (*model.CreateOrderResp, error) {
	if rspMap.GetString("success") == SUCCESS_FALSE {
		return nil, ierr.NewIError(ierr.ParamErr, fmt.Sprintf("API错误: %s(%s)", rspMap.Get("errorMsg"), rspMap.Get("errorCode")))
	}

	dataRes := rspMap.GetMap("data")
	result := model.CreateOrderResp{
		OrderId:   dataRes.GetString("orderNo"),
		WaybillNo: dataRes.GetString("waybillNo"),
	}

	return &result, nil
}

type CancelOrderAdaptor struct {
}

func (c *CancelOrderAdaptor) ConvertRequest(cfg *Config, req *model.CancelOrderReq) *CancelOrderRequest {
	return &CancelOrderRequest{
		BillCode:    req.WaybillNo,
		OrderType:   "02",
		OrderSource: cfg.SourceCode,
	}
}

func (c *CancelOrderAdaptor) ParseResponse(rspMap goutil.Map) error {
	if rspMap.GetString("success") == SUCCESS_FALSE {
		return ierr.NewIError(ierr.ParamErr, fmt.Sprintf("API错误: %s(%s)", rspMap.Get("errorMsg"), rspMap.Get("errorCode")))
	}

	return nil
}

type GetPriceQuoteReqAdaptor struct {
}

func (c *GetPriceQuoteReqAdaptor) ConvertRequest(req *model.GetPriceQuoteReq) *model.GetPriceQuoteReq {
	return req
}

func (c *GetPriceQuoteReqAdaptor) ParseResponse(rspMap goutil.Map) (*model.PriceQuote, error) {
	if rspMap.GetString("success") == SUCCESS_FALSE {
		return nil, ierr.NewIError(ierr.ParamErr, fmt.Sprintf("API错误: %s(%s); data: %v", rspMap.Get("errorMsg"), rspMap.Get("errorCode"), rspMap))
	}

	dataRes := rspMap.GetMapArrayP("data/AvailableServiceItemList")
	if len(dataRes) == 0 {
		return nil, ierr.NewIError(ierr.InternalError, fmt.Sprintf("data: %v", rspMap))
	}
	feeModel := dataRes[0].GetMap("feeModel")
	var result = model.PriceQuote{
		StartPrice:          stringToInt(feeModel.GetString("startPrice")),
		ContinuedHeavy:      stringToInt(feeModel.GetString("continuedHeavy")),
		StartWeight:         stringToInt(feeModel.GetString("startWeight")),
		ContinuedHeavyPrice: stringToInt(feeModel.GetString("continuedHeavyPrice")),
		TotalPrice:          stringToInt(feeModel.GetString("totalPrice")),
	}

	return &result, nil
}

func stringToInt(numStr string) int {
	startPrice, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}
	return startPrice
}

type SubscribeTrackingAdaptor struct {
}

func (c *SubscribeTrackingAdaptor) ConvertRequest(req *model.SubscribeTrackingReq) goutil.Map {
	return goutil.Map{
		"waybillNo": req.WaybillNo,
	}
}

func (c *SubscribeTrackingAdaptor) ParseResponse(rspMap goutil.Map) error {
	if rspMap.GetString("success") == SUCCESS_FALSE {
		return ierr.NewIError(ierr.ParamErr, fmt.Sprintf("API错误: %s(%s); data: %v", rspMap.Get("errorMsg"), rspMap.Get("errorCode"), rspMap))
	}
	return nil
}

// QueryLogisticsAdaptor 物流查询适配器
type QueryLogisticsAdaptor struct {
}

func (q *QueryLogisticsAdaptor) ConvertRequest(req *model.QueryLogisticsRequest) goutil.Map {
	data := goutil.Map{
		"order":         "desc",
		"waybillNoList": []string{req.WaybillNo},
	}
	return data
}

func (q *QueryLogisticsAdaptor) ParseResponse(waybillNo string, rspMap goutil.Map) (*model.QueryResp, error) {
	if rspMap.GetString("success") == SUCCESS_FALSE {
		return nil, ierr.NewIError(ierr.ParamErr, fmt.Sprintf("API错误: %s(%s)", rspMap.Get("errorMsg"), rspMap.Get("errorCode")))
	}
	dataRes := rspMap.GetMapArrayP("data/" + waybillNo)
	if len(dataRes) == 0 {
		return nil, ierr.NewIError(ierr.ParamErr, fmt.Sprintf("数据解析错误: %v", rspMap))
	}
	var data []*model.Data
	for i := len(dataRes) - 1; i >= 0; i-- {
		m := dataRes[i]
		data = append(data, &model.Data{
			Time:       m.GetString("opTime"),
			Context:    m.GetString("memo"),
			Ftime:      m.GetString("opTime"),
			AreaCode:   "",
			AreaName:   m.GetString("opOrgProvinceName") + "," + m.GetString("opOrgCityName"),
			Status:     m.GetString("status"),
			Location:   m.GetString("opOrgCityName"),
			AreaCenter: "",
			AreaPinYin: "",
			StatusCode: "",
		})
	}

	//for _, m := range dataRes {

	//}
	//
	//result := model.QueryResp{
	//	WaybillNo:   rspMap.GetString("nu"),
	//	Ischeck:     rspMap.GetString("ischeck"),
	//	ExpressCode: rspMap.GetString("com"),
	//	//Status:      GetHighLevelStatusMeaning(statusCode),
	//	State: rspMap.GetString("state"),
	//	Data:  data,
	//}

	return nil, nil
}
