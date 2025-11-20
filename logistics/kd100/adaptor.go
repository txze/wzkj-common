package kd100

import (
	"strconv"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/logistics/model"
	"github.com/txze/wzkj-common/pkg/util"
)

// HighLevelStatusMeaning 为快递100高级物流状态含义的映射表
var HighLevelStatusMeaning = map[int]string{
	// 揽收
	1:   "揽件",
	101: "已经下快件单",
	102: "待快递公司揽收",
	103: "快递公司已经揽收",

	// 在途
	0:    "在途中",
	1001: "快件到达收件人城市",
	1002: "快件处于运输过程中",
	1003: "快件发往到新的收件地址",

	// 派件
	5:   "正在派件",
	501: "快件已经投递到快递柜或者快递驿站",

	// 签收
	3:   "已签收",
	301: "收件人正常签收",
	302: "快件显示派件异常，但后续正常签收",
	303: "快件已被代签",
	304: "快件已从快递柜或者驿站取出签收",

	// 退回
	6: "退回",

	// 退签
	4:   "此快件单已退签",
	401: "此快件单已撤销",

	// 拒签
	14: "收件人拒签快件",

	// 转投
	7: "快件转给其他快递公司邮寄",

	// 疑难
	2:   "快件存在疑难",
	201: "快件长时间派件后未签收",
	202: "快件长时间没有派件或签收",
	203: "收件人发起拒收快递，待发货方确认",
	204: "快件派件时遇到异常情况",
	205: "快件在快递柜或者驿站长时间未取",
	206: "无法联系到收件人",
	207: "超出快递公司的服务区范围",
	208: "快件滞留在网点，没有派送",
	209: "快件破损",
	210: "寄件人申请撤销寄件",

	// 清关
	8:  "快件清关",
	10: "快件等待清关",
	11: "快件正在清关流程中",
	12: "快件已完成清关流程",
	13: "货物在清关过程中出现异常",
}

// GetHighLevelStatusMeaning 获取高级物流状态对应的含义
func GetHighLevelStatusMeaning(status int) string {
	if v, ok := HighLevelStatusMeaning[status]; ok {
		return v
	}
	return "未知状态"
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

	result := model.QueryResp{
		WaybillNo:   rspMap.GetString("nu"),
		Ischeck:     rspMap.GetString("ischeck"),
		ExpressCode: rspMap.GetString("com"),
		Status:      GetHighLevelStatusMeaning(statusCode),
		State:       rspMap.GetString("state"),
		Data:        data,
	}

	return &result, nil
}
