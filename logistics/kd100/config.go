package kd100

const (
	QueryURL        = "https://poll.kuaidi100.com/poll/query.do"
	ParseAddressURL = "https://api.kuaidi100.com/address/resolution"
)

type KD100Config struct {
	KEY      string
	CUSTOMER string
	Secret   string
	Resultv2 string //获取高级状态名称及状态值 传 “4” 或者“8”
}
