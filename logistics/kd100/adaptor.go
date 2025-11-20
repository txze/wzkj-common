package kd100

import (
	"strconv"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/pkg/util"
)

// StatusGroup 每个物流状态分组
type StatusGroup struct {
	GroupName string
	Items     map[int]string
}

// HighLevelStatusGroups 按分组整理后的高级物流状态
var HighLevelStatusGroups = []StatusGroup{
	{
		GroupName: "揽收",
		Items: map[int]string{
			1:   "揽件",
			101: "已下单",
			102: "待揽收",
			103: "已揽收",
		},
	},
	{
		GroupName: "在途",
		Items: map[int]string{
			0:    "在途",
			1001: "到达派件城市",
			1002: "运输中",
			1003: "转递",
		},
	},
	{
		GroupName: "派件",
		Items: map[int]string{
			5:   "派件",
			501: "投柜或驿站",
		},
	},
	{
		GroupName: "签收",
		Items: map[int]string{
			3:   "已签收",
			301: "本人签收",
			302: "派件异常后签收",
			303: "已代签",
			304: "投柜或站签收",
		},
	},
	{
		GroupName: "退回",
		Items: map[int]string{
			6: "退回",
		},
	},
	{
		GroupName: "退签",
		Items: map[int]string{
			4:   "已退签",
			401: "已撤销",
			14:  "拒签",
		},
	},
	{
		GroupName: "转投",
		Items: map[int]string{
			7: "转投",
		},
	},
	{
		GroupName: "疑难",
		Items: map[int]string{
			2:   "快件存在疑难",
			201: "超时未签收",
			202: "超时未更新",
			203: "拒收",
			204: "派件异常",
			205: "柜或驿站超时未取",
			206: "无法联系到收件人",
			207: "超出服务范围",
			208: "滞留",
			209: "破损",
			210: "销单",
		},
	},
	{
		GroupName: "清关",
		Items: map[int]string{
			8:  "清关",
			10: "待清关",
			11: "清关中",
			12: "已清关",
			13: "清关异常",
		},
	},
}

// GetGroupByStatus 根据状态值获取分组名
func GetGroupByStatus(code int) (string, string) {
	for _, group := range HighLevelStatusGroups {
		if v, ok := group.Items[code]; ok {
			return group.GroupName, v
		}
	}
	return "未知分组", "未知状态"
}

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

	status, ItemStatus := GetGroupByStatus(statusCode)
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
