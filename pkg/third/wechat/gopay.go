package wechat

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat"

	"github.com/spf13/viper"
)

type PayClient struct {
	*wechat.Client

	RefundNotifyUrl string
}

var payClient *PayClient

func InitPay() error {
	if viper.GetBool("wechat.is_test") {
		return nil
	}
	payClient = new(PayClient)
	payClient.Client = wechat.NewClient(
		viper.GetString("wechat.appid"),
		viper.GetString("wechat.mchid"),
		viper.GetString("wechat.pay_key"),
		viper.GetBool("wechat.is_prod"),
	)
	payClient.RefundNotifyUrl = viper.GetString("wechat.refund_notify_url")
	var err error
	err = payClient.AddCertPemFilePath(
		viper.GetString("wechat.certFile"),
		viper.GetString("wechat.keyFile"),
	)
	if err != nil {
		return err
	}

	err = payClient.AddCertPkcs12FilePath(
		viper.GetString("wechat.pkcs12File"),
	)
	if err != nil {
		return err
	}
	return nil
}

func GetPayClient() *PayClient {
	return payClient
}

func (p *PayClient) DoRefund(bodyMap gopay.BodyMap) (wxRsp *wechat.RefundResponse, err error) {
	if _, ok := bodyMap["notify_url"]; !ok && p.RefundNotifyUrl != "" {
		bodyMap["notify_url"] = p.RefundNotifyUrl
	}
	wxRsp, _, err = payClient.Refund(context.Background(), bodyMap)
	if err != nil {
		return nil, err
	}
	if wxRsp.ReturnCode != gopay.SUCCESS {
		return nil, errors.New(wxRsp.ReturnMsg)
	}
	if wxRsp.ResultCode != gopay.SUCCESS {
		return nil, fmt.Errorf("code: %v, desc: %v", wxRsp.ErrCode, wxRsp.ErrCodeDes)
	}
	return wxRsp, nil
}

func (p *PayClient) DoTransfer(bodyMap gopay.BodyMap) (wxRsp *wechat.TransfersResponse, err error) {
	wxRsp, err = payClient.Transfer(context.Background(), bodyMap)
	if err != nil {
		return nil, err
	}
	if wxRsp.ReturnCode != gopay.SUCCESS {
		return nil, errors.New(wxRsp.ReturnMsg)
	}
	if wxRsp.ResultCode != gopay.SUCCESS {
		return nil, fmt.Errorf("code: %v, desc: %v", wxRsp.ErrCode, wxRsp.ErrCodeDes)
	}
	return wxRsp, nil
}
