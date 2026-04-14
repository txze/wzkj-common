package allinpay

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/txze/wzkj-common/pay/common"
)

func TestAllInPay_GetType(t *testing.T) {
	config := AllInPayConfig{}

	allinpay, err := NewAllInPay(config)
	if err != nil {
		t.Errorf("NewAllInPay returned error: %v", err)
	}

	request := &common.PaymentRequest{
		OrderNo:     fmt.Sprintf("TEST_%d", time.Now().Unix()),
		GoodsName:   "测试商品",
		ProductCode: "WX_SMALL_APP",
		Amount:      100,
		Expire:      "15",
	}

	ctx := context.Background()
	rsp, err := allinpay.Pay(ctx, request)
	if err != nil {
		t.Errorf("Pay returned error: %v", err)
	}
	if rsp == nil {
		t.Error("Pay returned nil response")
	}

	t.Log(rsp)
}
