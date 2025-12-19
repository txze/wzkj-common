package model

var StatusCodeMapping = map[string]int{
	"收件":    103,
	"发件":    1002,
	"到件":    1001,
	"派件":    5,
	"第三方代派": 5,
	"柜机代收":  501,
	"驿站代收":  501,
	"快件取出":  304,
	"驿站出库":  1003, // 推荐
	"问题件":   2,
	"退回件":   6,
	"改地址件":  7,
	"转邮":    7,
	"转同行":   7,
	"签收":    3,
	"客户签收":  3,
}

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

// FindGroupAndName 根据状态值获取分组名
func FindGroupAndName(code int) (string, string) {
	for _, group := range HighLevelStatusGroups {
		if v, ok := group.Items[code]; ok {
			return group.GroupName, v
		}
	}
	return "未知分组", "未知状态"
}

func ResolveStatusByText(status string) (groupName string, statusText string, code int) {
	code = StatusCodeMapping[status]
	groupName, statusText = FindGroupAndName(code)
	return
}
