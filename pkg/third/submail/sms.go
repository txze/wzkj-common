package submail

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/hzxiao/goutil"
)

type Config struct {
	AppID     string `mapstructure:"appid"`
	Project   string `mapstructure:"project"`
	Signature string `mapstructure:"signature"`
}

const SMS_STATUS_SUCCESS = "success"
const SMS_STATUS_ERROR = "error"

const SMS_URI_JSON = "https://api-v4.mysubmail.com/sms/xsend.json"

var globalConfig *Config

func InitSMS(cfg *Config) {
	globalConfig = cfg
}

func SendSMS(to string, vars string) error {
	if globalConfig == nil {
		return fmt.Errorf("SMS未初始化，请先调用InitSMS()")
	}
	return sendSMS(globalConfig.AppID, globalConfig.Project, to, vars, globalConfig.Signature)
}

func SendSMSWithConfig(cfg *Config, to string, vars string) error {
	return sendSMS(cfg.AppID, cfg.Project, to, vars, cfg.Signature)
}

func sendSMS(appid, project string, to, vars, signature string) error {
	var result goutil.Map
	var client = resty.New()

	var paylod = map[string]string{
		"appid":     appid,
		"to":        to,
		"project":   project,
		"vars":      vars,
		"signature": signature,
	}
	resp, err := client.R().
		SetBody(paylod).
		SetResult(&result).
		Post(SMS_URI_JSON)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http error, code: %d, error: %s", resp.StatusCode(), resp.String())
	}

	if result.GetString("status") != SMS_STATUS_SUCCESS {
		return fmt.Errorf("code: %d, error: %s", result.GetInt64("code"), result.GetString("msg"))
	}

	return nil
}
