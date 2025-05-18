package alipay

import (
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
)

type AliPayClient struct {
	*alipay.Client
}

var aliPayClient *AliPayClient
var alipayAppClient *AliPayClient

func InitAliPay() error {
	var client, err = alipay.New(
		viper.GetString("alipay.appid"),
		viper.GetString("alipay.private_key"),
		viper.GetBool("alipay.is_prod"),
		alipay.WithProductionGateway(""))
	if err != nil {
		return err
	}

	// 加载应用公钥证书
	if err = client.LoadAppCertPublicKeyFromFile(viper.GetString("alipay.app_public_cert")); err != nil {
		// 错误处理
		return err
	}

	// 加载支付宝根证书
	if err = client.LoadAliPayRootCertFromFile(viper.GetString("alipay.alipay_root_cert")); err != nil {
		// 错误处理
		return err
	}

	// 加载支付宝公钥证书
	if err = client.LoadAlipayCertPublicKeyFromFile(viper.GetString("alipay.alipay_public_cert")); err != nil {
		// 错误处理
		return err
	}

	// 加载内容密钥，可选
	var encryptKey = viper.GetString("alipay.encrypt_key")
	if encryptKey != "" {
		if err = client.SetEncryptKey(encryptKey); err != nil {
			// 错误处理
			return err
		}
	}

	aliPayClient = new(AliPayClient)
	aliPayClient.Client = client

	return nil
}

func InitAliAppPay() error {
	var client, err = alipay.New(
		viper.GetString("alipay_app.appid"),
		viper.GetString("alipay_app.private_key"),
		viper.GetBool("alipay_app.is_prod"),
		alipay.WithProductionGateway(""))
	if err != nil {
		return err
	}

	// 加载应用公钥证书
	if err = client.LoadAppCertPublicKeyFromFile(viper.GetString("alipay_app.app_public_cert")); err != nil {
		// 错误处理
		return err
	}

	// 加载支付宝根证书
	if err = client.LoadAliPayRootCertFromFile(viper.GetString("alipay_app.alipay_root_cert")); err != nil {
		// 错误处理
		return err
	}

	// 加载支付宝公钥证书
	if err = client.LoadAlipayCertPublicKeyFromFile(viper.GetString("alipay_app.alipay_public_cert")); err != nil {
		// 错误处理
		return err
	}

	// 加载内容密钥，可选
	if err = client.SetEncryptKey(viper.GetString("alipay_app.encrypt_key")); err != nil {
		// 错误处理
		return err
	}

	alipayAppClient = new(AliPayClient)
	alipayAppClient.Client = client

	return nil
}

func GetAliPayClient() *AliPayClient {
	return aliPayClient
}

func GetAliAppPayClient() *AliPayClient {
	return alipayAppClient
}
