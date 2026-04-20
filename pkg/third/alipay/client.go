package alipay

import (
	"github.com/smartwalle/alipay/v3"
)

type Config struct {
	AppID            string `mapstructure:"appid"`
	PrivateKey       string `mapstructure:"private_key"`
	IsProd           bool   `mapstructure:"is_prod"`
	AppPublicCert    string `mapstructure:"app_public_cert"`
	AlipayRootCert   string `mapstructure:"alipay_root_cert"`
	AlipayPublicCert string `mapstructure:"alipay_public_cert"`
	EncryptKey       string `mapstructure:"encrypt_key"`
}

type AliPayClient struct {
	*alipay.Client
}

var aliPayClient *AliPayClient
var alipayAppClient *AliPayClient

func InitAliPayWithConfig(cfg *Config) error {
	var client, err = alipay.New(
		cfg.AppID,
		cfg.PrivateKey,
		cfg.IsProd,
		alipay.WithProductionGateway(""))
	if err != nil {
		return err
	}

	if err = client.LoadAppCertPublicKeyFromFile(cfg.AppPublicCert); err != nil {
		return err
	}

	if err = client.LoadAliPayRootCertFromFile(cfg.AlipayRootCert); err != nil {
		return err
	}

	if err = client.LoadAlipayCertPublicKeyFromFile(cfg.AlipayPublicCert); err != nil {
		return err
	}

	var encryptKey = cfg.EncryptKey
	if encryptKey != "" {
		if err = client.SetEncryptKey(encryptKey); err != nil {
			return err
		}
	}

	aliPayClient = new(AliPayClient)
	aliPayClient.Client = client

	return nil
}

func InitAliAppPayWithConfig(cfg *Config) error {
	var client, err = alipay.New(
		cfg.AppID,
		cfg.PrivateKey,
		cfg.IsProd,
		alipay.WithProductionGateway(""))
	if err != nil {
		return err
	}

	if err = client.LoadAppCertPublicKeyFromFile(cfg.AppPublicCert); err != nil {
		return err
	}

	if err = client.LoadAliPayRootCertFromFile(cfg.AlipayRootCert); err != nil {
		return err
	}

	if err = client.LoadAlipayCertPublicKeyFromFile(cfg.AlipayPublicCert); err != nil {
		return err
	}

	if err = client.SetEncryptKey(cfg.EncryptKey); err != nil {
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
