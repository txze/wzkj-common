package logistics

import "fmt"

// LogisticsStrategy 定义接口策略
type LogisticsStrategy interface {
	// QueryLogisticsByNumber 查询物流信息
	QueryLogisticsByNumber(string, string, string, string) (string, error)
}

type LogisticsResponse struct {
	Message   string `json:"message"`
	Nu        string `json:"nu"`
	Ischeck   string `json:"ischeck"`
	Condition string `json:"condition"`
	Com       string `json:"com"`
	Status    string `json:"status"`
	State     string `json:"state"`
	Data      []struct {
		Time    string `json:"time"`
		Ftime   string `json:"ftime"`
		Context string `json:"context"`
	} `json:"data"`
}

// 物流策略
type LogisticsContext struct {
	strategy LogisticsStrategy
}

func (c *LogisticsContext) SetStrategy(strategy LogisticsStrategy) {
	c.strategy = strategy
}

func (c *LogisticsContext) QueryLogisticsByNumber(code, number, phone, resultv2 string) (string, error) {
	if c.strategy == nil {
		return "", fmt.Errorf("未设置物流策略")
	}
	return c.strategy.QueryLogisticsByNumber(code, number, phone, resultv2)
}
