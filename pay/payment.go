package pay

// PaymentStrategy 支付策略接口
// RQ 参数
// RS 返回
type PaymentStrategy[RQ, Notify, RS any] interface {
	// Pay 发起支付
	Pay(request *RQ) (map[string]interface{}, error)

	// VerifyNotification 验证支付通知
	VerifyNotification(notification *Notify) (bool, error)

	// QueryPayment 查询支付状态
	QueryPayment(orderID string) (*RS, error)

	// Refund 退款
	Refund(orderID string, amount float64) error

	// GenerateSign 生成签名
	GenerateSign(params map[string]interface{}) (string, error)

	// VerifySign 验证签名
	VerifySign(params map[string]interface{}, sign string) (bool, error)

	// GetType 获取支付类型
	GetType() string
}

type Payment[RQ, Notify, RS any] struct {
}

func NewPayment[RQ, Notify, RS any]() *Payment[RQ, Notify, RS] {
	return &Payment[RQ, Notify, RS]{}
}

func (p *Payment[RQ, Notify, RS]) SetStrategy(strategy PaymentStrategy[RQ, Notify, RS]) PaymentStrategy[RQ, Notify, RS] {
	if strategy == nil {
		return nil
	}

	return strategy
}

//
//func (p *Payment[RQ, RS]) Process(req RQ) (*RT, error) {
//	if p.strategy == nil {
//		return nil, nil
//	}
//	return p.strategy.Process(params)
//}
//
//func (p *Payment[CT, PT, RT]) Notify(data string) (bool, error) {
//	if p.strategy == nil {
//		return false, nil
//	}
//	return p.strategy.Notify(data)
//}
