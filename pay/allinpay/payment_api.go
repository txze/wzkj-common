package allinpay

import (
	"context"
	"strconv"

	"github.com/hzxiao/goutil"

	"github.com/txze/wzkj-common/pay/common"
)

// WxSmallPayStrategy 微信小程序支付策略
type WxSmallPayStrategy struct{}

// NewWxSmallPayStrategy 创建微信小程序支付策略
func NewWxSmallPayStrategy() *WxSmallPayStrategy {
	return &WxSmallPayStrategy{}
}

// Process 处理微信小程序支付请求
func (s *WxSmallPayStrategy) Process(ctx context.Context, request *common.PaymentRequest) (goutil.Map, error) {
	params := goutil.Map{}
	params["goodsName"] = request.GoodsName
	params["outOrderNo"] = request.OrderNo
	params["transAmt"] = request.Amount
	if request.ProductCode == "WX_SMALL_APP" {
		params["payType"] = PayTypeApp
	} else if request.ProductCode == "WX_SMALL" {
		params["payType"] = PayTypeSmall
	}
	params["paySource"] = PaySourceDefault
	if request.Expire != "" {
		params["timeoutExpress"], _ = strconv.Atoi(request.Expire)
	}

	return params, nil
}

// GetUrl 获取支付接口路径
func (s *WxSmallPayStrategy) GetUrl() string {
	return WxSmallUrl
}
