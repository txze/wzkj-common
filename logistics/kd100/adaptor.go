package kd100

import (
	"strconv"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/pkg/ierr"
	"github.com/txze/wzkj-common/pkg/util"
)

// QueryLogisticsAdaptor 物流查询适配器
type QueryLogisticsAdaptor struct {
}

func (q *QueryLogisticsAdaptor) ConvertRequest(req *model.QueryLogisticsRequest) string {
	data := goutil.Map{
		"com":      req.ExpressCode,
		"num":      req.WaybillNo,
		"phone":    req.Phone,
		"resultv2": "4",
	}
	return util.S2Json(data)
}

func (q *QueryLogisticsAdaptor) ParseResponse(rspMap goutil.Map) (*model.QueryResp, error) {
	if rspMap.GetBool("result") == false && rspMap.GetString("message") != "ok" {
		return nil, ierr.NewIErrorf(ierr.InternalError, "error:%s code:%s", rspMap.GetString("message"), rspMap.GetString("returnCode"))
	}
	dataRes := rspMap.GetMapArray("data")
	var statusCode int
	var err error
	if len(dataRes) > 0 {
		statusCodeStr := dataRes[0].GetString("statusCode")
		statusCode, err = strconv.Atoi(statusCodeStr)
		if err != nil {
			return nil, err
		}
	}
	var data []*model.Data
	for _, m := range dataRes {
		data = append(data, &model.Data{
			Time:       m.GetString("time"),
			Context:    m.GetString("context"),
			Ftime:      m.GetString("ftime"),
			AreaCode:   m.GetString("areaCode"),
			AreaName:   m.GetString("areaName"),
			Status:     m.GetString("status"),
			Location:   m.GetString("location"),
			AreaCenter: m.GetString("areaCenter"),
			AreaPinYin: m.GetString("areaPinYin"),
			StatusCode: m.GetString("statusCode"),
		})
	}

	status, ItemStatus := model.FindGroupAndName(statusCode)
	result := model.QueryResp{
		WaybillNo:   rspMap.GetString("nu"),
		Ischeck:     rspMap.GetString("ischeck"),
		ExpressCode: rspMap.GetString("com"),
		Status:      status,
		ItemStatus:  ItemStatus,
		State:       rspMap.GetString("state"),
		Data:        data,
	}

	return &result, nil
}
