package pay

// PaymentStrategy 支付策略接口
type PaymentStrategy[CT, PT, RT any] interface {
	SetConfig(CT)
	Process(PT) (*RT, error)
	Notify(data string) (bool, error)
}

type Payment[CT, PT, RT any] struct {
	strategy PaymentStrategy[CT, PT, RT]
}

func NewPayment[CT, PT, RT any]() *Payment[CT, PT, RT] {
	return &Payment[CT, PT, RT]{}
}

func (p *Payment[CT, PT, RT]) SetStrategy(strategy PaymentStrategy[CT, PT, RT]) {
	if strategy == nil {
		return
	}
	p.strategy = strategy
}

func (p *Payment[CT, PT, RT]) SetConfig(config CT) {
	if p.strategy == nil {
		return
	}
	p.strategy.SetConfig(config)
}

func (p *Payment[CT, PT, RT]) Process(params PT) (*RT, error) {
	if p.strategy == nil {
		return nil, nil
	}
	return p.strategy.Process(params)
}

func (p *Payment[CT, PT, RT]) Notify(data string) (bool, error) {
	if p.strategy == nil {
		return false, nil
	}
	return p.strategy.Notify(data)
}
