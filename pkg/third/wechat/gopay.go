package wechat

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat"
)

type Config struct {
	IsTest          bool   `mapstructure:"is_test"`
	AppID           string `mapstructure:"appid"`
	MchID           string `mapstructure:"mchid"`
	PayKey          string `mapstructure:"pay_key"`
	IsProd          bool   `mapstructure:"is_prod"`
	RefundNotifyUrl string `mapstructure:"refund_notify_url"`
	CertFile        string `mapstructure:"certFile"`
	KeyFile         string `mapstructure:"keyFile"`
	Pkcs12File      string `mapstructure:"pkcs12File"`
}

type PayClient struct {
	*wechat.Client

	RefundNotifyUrl string
}

var payClient *PayClient

func InitPayWithConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("配置不能为空")
	}

	if cfg.IsTest {
		return nil
	}

	payClient = new(PayClient)
	payClient.Client = wechat.NewClient(
		cfg.AppID,
		cfg.MchID,
		cfg.PayKey,
		cfg.IsProd,
	)
	payClient.RefundNotifyUrl = cfg.RefundNotifyUrl

	var err error
	err = payClient.AddCertPemFilePath(
		cfg.CertFile,
		cfg.KeyFile,
	)
	if err != nil {
		return err
	}

	err = payClient.AddCertPkcs12FilePath(
		cfg.Pkcs12File,
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
