package logistics

import (
	"errors"

	"github.com/txze/wzkj-common/logistics/model"
)

// 物流策略
type LogisticsContext struct {
	strategy LogisticsProvider
}

func (c *LogisticsContext) SetStrategy(strategy LogisticsProvider) {
	c.strategy = strategy
}

func (c *LogisticsContext) QueryLogisticsByNumber(req *model.QueryLogisticsRequest) (*model.QueryResp, error) {
	if c.strategy == nil {
		return nil, errors.New("logistics strategy is nil")
	}
	return c.strategy.QueryLogistics(req)
}

func (c *LogisticsContext) ParseAddress(addr string) (model.Address, error) {
	if c.strategy == nil {
		return nil, errors.New("logistics strategy is nil")
	}
	return c.strategy.ParseAddress(addr)
}
