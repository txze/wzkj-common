package wechat

import (
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
)

type MiniProgramConfig struct {
	AppID        string `mapstructure:"appid"`
	AppSecret    string `mapstructure:"app_secret"`
	MchID        string `mapstructure:"mchid"`
	PayNotifyURL string `mapstructure:"pay_notify_url"`
	PayKey       string `mapstructure:"pay_key"`
}

type H5Config struct {
	AppID        string `mapstructure:"h5_appid"`
	AppSecret    string `mapstructure:"h5_app_secret"`
	MchID        string `mapstructure:"mchid"`
	PayNotifyURL string `mapstructure:"pay_notify_url"`
	PayKey       string `mapstructure:"pay_key"`
}

func InitWxMiniProgramWithConfig(cfg *MiniProgramConfig) {
	memcache := cache.NewMemory()
	var wxcfg = &wechat.Config{
		AppID:          cfg.AppID,
		AppSecret:      cfg.AppSecret,
		Token:          "",
		EncodingAESKey: "",
		PayMchID:       cfg.MchID,
		PayNotifyURL:   cfg.PayNotifyURL,
		PayKey:         cfg.PayKey,
		Cache:          memcache,
	}

	wxMiniprogramClient = &WxClient{
		Wechat: wechat.NewWechat(wxcfg),
		Oauth:  &Oauth{},
	}
	wxMiniprogramClient.Oauth.Oauth = wxMiniprogramClient.Wechat.GetOauth()
}

func InitSlaveWxMiniProgramWithConfig(cfg *MiniProgramConfig) {
	memcache := cache.NewMemory()
	var wxcfg = &wechat.Config{
		AppID:          cfg.AppID,
		AppSecret:      cfg.AppSecret,
		Token:          "",
		EncodingAESKey: "",
		PayMchID:       cfg.MchID,
		PayNotifyURL:   cfg.PayNotifyURL,
		PayKey:         cfg.PayKey,
		Cache:          memcache,
	}

	wxSlaveMiniprogramClient = &WxClient{
		Wechat: wechat.NewWechat(wxcfg),
		Oauth:  &Oauth{},
	}
	wxSlaveMiniprogramClient.Oauth.Oauth = wxSlaveMiniprogramClient.Wechat.GetOauth()
}

func InitWxH5WithConfig(cfg *H5Config) {
	memcache := cache.NewMemory()
	var wxcfg = &wechat.Config{
		AppID:          cfg.AppID,
		AppSecret:      cfg.AppSecret,
		Token:          "",
		EncodingAESKey: "",
		PayMchID:       cfg.MchID,
		PayNotifyURL:   cfg.PayNotifyURL,
		PayKey:         cfg.PayKey,
		Cache:          memcache,
	}

	wxH5Client = &WxClient{
		Wechat: wechat.NewWechat(wxcfg),
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

func GetClient() *WxClient {
	return wxMiniprogramClient
}

func GetSlaveClient() *WxClient {
	return wxSlaveMiniprogramClient
}

func GetWxH5Client() *WxClient {
	return wxH5Client
}
