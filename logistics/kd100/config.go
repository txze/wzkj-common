package kd100

const (
	QueryURL        = "https://poll.kuaidi100.com/poll/query.do"
	ParseAddressURL = "https://api.kuaidi100.com/address/resolution"
)

type KD100Config struct {
	KEY      string
	CUSTOMER string
	Secret   string
}
