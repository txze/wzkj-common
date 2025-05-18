package wechat

import (
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"github.com/spf13/viper"
)

// 初始化小程序
func InitWxMiniProgram() {
	memcache := cache.NewMemory()
	var cfg = &wechat.Config{
		AppID:          viper.GetString("wechat.appid"),
		AppSecret:      viper.GetString("wechat.app_secret"),
		Token:          "",
		EncodingAESKey: "",
		PayMchID:       viper.GetString("wechat.mchid"),
		PayNotifyURL:   viper.GetString("wechat.pay_notify_url"),
		PayKey:         viper.GetString("wechat.pay_key"),
		Cache:          memcache,
	}

	// 初始化默认wechat client对象
	wxMiniprogramClient = &WxClient{
		Wechat: wechat.NewWechat(cfg),
		Oauth:  &Oauth{},
	}
	wxMiniprogramClient.Oauth.Oauth = wxMiniprogramClient.Wechat.GetOauth()
}

// 初始化小程序
func InitSlaveWxMiniProgram() {
	memcache := cache.NewMemory()
	var cfg = &wechat.Config{
		AppID:          viper.GetString("wechat.slave_appid"),
		AppSecret:      viper.GetString("wechat.slave_app_secret"),
		Token:          "",
		EncodingAESKey: "",
		PayMchID:       viper.GetString("wechat.mchid"),
		PayNotifyURL:   viper.GetString("wechat.pay_notify_url"),
		PayKey:         viper.GetString("wechat.pay_key"),
		Cache:          memcache,
	}

	// 初始化默认wechat client对象
	wxSlaveMiniprogramClient = &WxClient{
		Wechat: wechat.NewWechat(cfg),
		Oauth:  &Oauth{},
	}
	wxSlaveMiniprogramClient.Oauth.Oauth = wxSlaveMiniprogramClient.Wechat.GetOauth()
}

func InitWxH5() {
	memcache := cache.NewMemory()
	var cfg = &wechat.Config{
		AppID:          viper.GetString("wechat.h5_appid"),
		AppSecret:      viper.GetString("wechat.h5_app_secret"),
		Token:          "",
		EncodingAESKey: "",
		PayMchID:       viper.GetString("wechat.mchid"),
		PayNotifyURL:   viper.GetString("wechat.pay_notify_url"),
		PayKey:         viper.GetString("wechat.pay_key"),
		Cache:          memcache,
	}

	// 初始化默认wechat client对象
	wxH5Client = &WxClient{
		Wechat: wechat.NewWechat(cfg),
		Oauth:  &Oauth{},
	}
	wxH5Client.Oauth.Oauth = wxH5Client.Wechat.GetOauth()
}

var wxMiniprogramClient *WxClient
var wxSlaveMiniprogramClient *WxClient
var wxH5Client *WxClient

type WxClient struct {
	*wechat.Wechat
	Oauth *Oauth
}

/*
默认是微信小程序
*/
func GetClient() *WxClient {
	return wxMiniprogramClient
}

/*
支付子微信小程序
*/
func GetSlaveClient() *WxClient {
	return wxSlaveMiniprogramClient
}

func GetWxH5Client() *WxClient {
	return wxH5Client
}
