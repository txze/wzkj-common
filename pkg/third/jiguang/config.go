package jiguang

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 极光一键登录配置
type Config struct {
	AppKey       string `mapstructure:"app_key"`     // 应用Key
	MasterSecret string `mapstructure:"secret"`      // 主密钥
	PrivateKey   string `mapstructure:"private_key"` // RSA私钥
	APIURL       string `mapstructure:"api_url"`     // API地址
}

// DefaultConfig 获取默认配置
func DefaultConfig() *Config {
	return &Config{
		APIURL: "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
	}
}

// LoadConfig 从viper加载配置
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// 从viper获取配置
	appKey := viper.GetString("jiguang.app_key")
	masterSecret := viper.GetString("jiguang.secret")
	privateKey := viper.GetString("jiguang.private_key")
	apiURL := viper.GetString("jiguang.api_url")

	// 验证必填配置
	if appKey == "" {
		return nil, fmt.Errorf("极光配置缺失: jiguang.app_key")
	}
	if masterSecret == "" {
		return nil, fmt.Errorf("极光配置缺失: jiguang.master_secret")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("极光配置缺失: jiguang.private_key")
	}

	config.AppKey = appKey
	config.MasterSecret = masterSecret
	config.PrivateKey = privateKey

	// 如果配置了自定义API地址，使用自定义地址
	if apiURL != "" {
		config.APIURL = apiURL
	}

	return config, nil
}

// Validate 验证配置
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

	// 验证私钥格式
	if !strings.Contains(c.PrivateKey, "-----BEGIN") {
		return fmt.Errorf("PrivateKey格式不正确，应该包含PEM格式的RSA私钥")
	}

	return nil
}
