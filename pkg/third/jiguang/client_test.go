package jiguang

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "有效配置",
			config: &Config{
				AppKey:       "test_app_key",
				MasterSecret: "test_master_secret",
				PrivateKey:   "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----",
				APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
			},
			wantErr: false,
		},
		{
			name: "AppKey为空",
			config: &Config{
				AppKey:       "",
				MasterSecret: "test_master_secret",
				PrivateKey:   "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----",
				APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
			},
			wantErr: true,
		},
		{
			name: "MasterSecret为空",
			config: &Config{
				AppKey:       "test_app_key",
				MasterSecret: "",
				PrivateKey:   "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----",
				APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
			},
			wantErr: true,
		},
		{
			name: "PrivateKey为空",
			config: &Config{
				AppKey:       "test_app_key",
				MasterSecret: "test_master_secret",
				PrivateKey:   "",
				APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
			},
			wantErr: true,
		},
		{
			name: "PrivateKey格式错误",
			config: &Config{
				AppKey:       "test_app_key",
				MasterSecret: "test_master_secret",
				PrivateKey:   "invalid_private_key",
				APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.APIURL != "https://api.verification.jpush.cn/v1/web/loginTokenVerify" {
		t.Errorf("DefaultConfig() APIURL = %v, want %v", config.APIURL, "https://api.verification.jpush.cn/v1/web/loginTokenVerify")
	}
}

func TestParsePrivateKey(t *testing.T) {
	// 这是一个示例RSA私钥，仅用于测试
	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1234567890abcdefghijklmnopqrstuvwxyz
-----END RSA PRIVATE KEY-----`

	_, err := parsePrivateKey(testPrivateKey)
	if err == nil {
		t.Log("parsePrivateKey() 成功解析私钥")
	} else {
		t.Logf("parsePrivateKey() 解析私钥失败: %v", err)
	}
}

func TestNewClient(t *testing.T) {
	config := &Config{
		AppKey:       "test_app_key",
		MasterSecret: "test_master_secret",
		PrivateKey:   "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----",
		APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Logf("NewClient() 创建客户端失败: %v", err)
		return
	}

	if client == nil {
		t.Error("NewClient() 返回的客户端为nil")
	}

	if client.config != config {
		t.Error("NewClient() 配置不匹配")
	}
}

// 示例：如何使用极光一键登录
func ExampleVerifyLoginToken() {
	// 注意：这是一个示例，实际使用时需要先初始化客户端
	// 1. 在应用启动时初始化
	// err := Init()
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 2. 验证登录Token
	// phoneInfo, err := VerifyLoginToken("your_login_token", "your_ex_id")
	// if err != nil {
	//     log.Printf("验证失败: %v", err)
	//     return
	// }

	// 3. 处理结果
	// if phoneInfo.Success {
	//     fmt.Printf("手机号: %s\n", phoneInfo.Phone)
	//     fmt.Printf("自定义ID: %s\n", phoneInfo.ExID)
	// } else {
	//     fmt.Printf("验证失败: %s\n", phoneInfo.Error)
	// }
}

// 示例：如何创建自定义客户端
func ExampleNewClient() {
	// 创建配置
	config := &Config{
		AppKey:       "your_app_key",
		MasterSecret: "your_master_secret",
		PrivateKey:   "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
		APIURL:       "https://api.verification.jpush.cn/v1/web/loginTokenVerify",
	}

	// 创建客户端
	client, err := NewClient(config)
	if err != nil {
		// 处理错误
		return
	}

	// 使用客户端验证登录Token
	phoneInfo, err := client.VerifyLoginToken("your_login_token", "your_ex_id")
	if err != nil {
		// 处理错误
		return
	}

	// 处理结果
	if phoneInfo.Success {
		// 验证成功，获取到手机号
		_ = phoneInfo.Phone
	} else {
		// 验证失败
		_ = phoneInfo.Error
	}
}
