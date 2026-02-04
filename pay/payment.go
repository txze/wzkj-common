package pay

import (
	"context"
	"net/http"

	"github.com/hzxiao/goutil"

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
	QueryPayment(ctx context.Context, orderID string) (*common.UnifiedResponse, error)

	// Refund 退款
	Refund(ctx context.Context, request *common.RefundRequest) (*common.RefundOrderResponse, error)

	//查询退款详情
	QueryRefund(ctx context.Context, refundNo, orderNo string) (*common.RefundResponse, error)

	// GenerateSign 生成签名
	GenerateSign(params map[string]interface{}) (string, error)

	// VerifySign 验证签名
	VerifySign(params map[string]interface{}) (bool, error)

	Close(ctx context.Context, orderId string) (bool, error)

	// GetType 获取支付类型
	GetType() string

	// VerifySettleNotification 验证分账通知
	VerifySettleNotification(ctx context.Context, req *http.Request) (*common.SettleNotificationResponse, error)

	// TradeOrderSettle 交易分账
	TradeOrderSettle(ctx context.Context, request common.TradeRoyaltyRateQueryRequestInterface) (*common.TradeRoyaltyRateQueryResponse, error)

	// ConfirmSettle 结算确认（目前仅支付宝支持）
	ConfirmSettle(ctx context.Context, request common.SettleConfirmRequestInterface) (*common.SettleConfirmResponse, error)

	// MergePay 合并支付
	//MergePay(ctx context.Context, data gopay.BodyMap) (goutil.Map, error)
}

// 数据转换接口
type DataMapper interface {
	// MapToSettleConfirmRequest 映射分账确认请求参数
	MapToSettleConfirmRequest(data goutil.Map) common.SettleConfirmRequestInterface

	// MapToTradeRoyaltyRateQueryRequest 映射分账查询请求参数
	MapToTradeRoyaltyRateQueryRequest(data goutil.Map) common.TradeRoyaltyRateQueryRequestInterface
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
