package jiguang

import (
	"fmt"
	"strings"
)

type Config struct {
	AppKey       string `mapstructure:"app_key"`     // 应用Key
	MasterSecret string `mapstructure:"secret"`      // 主密钥
	PrivateKey   string `mapstructure:"private_key"` // RSA私钥
	APIURL       string `mapstructure:"api_url"`     // API地址
}

func DefaultConfig() *Config {
	return &Config{
		APIURL: "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
	}
}

func LoadConfigFromConfig(cfg *Config) (*Config, error) {
	if cfg == nil {
		return nil, fmt.Errorf("配置不能为空")
	}

	config := DefaultConfig()

	if cfg.AppKey != "" {
		config.AppKey = cfg.AppKey
	}
	if cfg.MasterSecret != "" {
		config.MasterSecret = cfg.MasterSecret
	}
	if cfg.PrivateKey != "" {
		config.PrivateKey = cfg.PrivateKey
	}
	if cfg.APIURL != "" {
		config.APIURL = cfg.APIURL
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.AppKey == "" {
		return fmt.Errorf("AppKey不能为空")
	}
	if c.MasterSecret == "" {
		return fmt.Errorf("MasterSecret不能为空")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("PrivateKey不能为空")
	}
	if c.APIURL == "" {
		return fmt.Errorf("APIURL不能为空")
	}

	if !strings.Contains(c.PrivateKey, "-----BEGIN") {
		return fmt.Errorf("PrivateKey格式不正确，应该包含PEM格式的RSA私钥")
	}

	return nil
}
