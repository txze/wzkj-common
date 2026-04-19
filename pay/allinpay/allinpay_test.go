package allinpay

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/txze/wzkj-common/pay/common"
)

func TestAllInPay_GetType(t *testing.T) {
	config := AllInPayConfig{
		AppId:      "wx0b06324c1324c33e",
		MchntId:    "510231",
		StoreId:    "510231",
		ChannelId:  "1055",
		PrivateKey: "MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgQfS7mPE9NYH7ItWbSVedRw4NzzZm3zljgvWEUCZf2vmgCgYIKoEcz1UBgi2hRANCAASUQKG8BdN94yB1OL1SKPiiWRcWahVL9QrTLFvFLANwS/w4FmGxyzfaMZHa852nJTpTKLDnZhc9gUDfRrC+eNaF",
		PublicKey:  "MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAElEChvAXTfeMgdTi9Uij4olkXFmoVS/UK0yxbxSwDcEv8OBZhscs32jGR2vOdpyU6Uyiw52YXPYFA30awvnjWhQ==",
		NotifyUrl:  "https://861nmir94998.vicp.fun/notify",
		IsProd:     true,
		SignType:   "SM2",
	}

	allinpay, err := NewAllInPay(config)
	if err != nil {
		t.Errorf("NewAllInPay returned error: %v", err)
	}

	request := &common.PaymentRequest{
		OrderNo:     fmt.Sprintf("%d", time.Now().UnixNano()),
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
