package pay

import (
	"context"
	"net/http"

	"github.com/txze/wzkj-common/pay/common"
)

// PaymentStrategy 支付策略接口
// RQ 参数
// RS 返回
type PaymentStrategy interface {
	// Pay 发起支付
	Pay(ctx context.Context, request *common.PaymentRequest) (map[string]interface{}, error)

	// VerifyNotification 验证支付通知
	VerifyNotification(req *http.Request) (*common.UnifiedResponse, error)

	// QueryPayment 查询支付状态
	QueryPayment(orderID string) (*common.UnifiedResponse, error)

	// Refund 退款
	Refund(orderID string, amount float64) error

	// GenerateSign 生成签名
	GenerateSign(params map[string]interface{}) (string, error)

	// VerifySign 验证签名
	VerifySign(params map[string]interface{}) (bool, error)

	Close(orderId string) (bool, error)

	// GetType 获取支付类型
	GetType() string
}

type Payment struct {
}

func NewPayment() *Payment {
	return &Payment{}
}

func (p *Payment) SetStrategy(strategy PaymentStrategy) PaymentStrategy {
	if strategy == nil {
		return nil
	}

	return strategy
}
