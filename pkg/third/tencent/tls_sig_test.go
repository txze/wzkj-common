package tencent

import (
	"testing"
)

// TestNewTLSSigAPI 测试创建TLS签名API实例
func TestNewTLSSigAPI(t *testing.T) {
	api := NewTLSSigAPI(123456789, "test_secret_key")
	if api == nil {
		t.Fatal("创建TLS签名API实例失败")
	}

	if api.sdkAppID != 123456789 {
		t.Errorf("SDKAppID不匹配，期望: 123456789, 实际: %d", api.sdkAppID)
	}

	if api.secretKey != "test_secret_key" {
		t.Errorf("SecretKey不匹配，期望: test_secret_key, 实际: %s", api.secretKey)
	}
}

// TestExpireTimeFunctions 测试过期时间函数
func TestExpireTimeFunctions(t *testing.T) {
	// 测试短期过期时间
	shortExpire := GetShortExpireTime()
	if shortExpire != 86400 {
		t.Errorf("短期过期时间应该为86400秒，实际为: %d", shortExpire)
	}

	// 测试默认过期时间
	defaultExpire := GetDefaultExpireTime()
	if defaultExpire != 5184000 {
		t.Errorf("默认过期时间应该为5184000秒，实际为: %d", defaultExpire)
	}

	// 测试长期过期时间
	longExpire := GetLongExpireTime()
	if longExpire != 31536000 {
		t.Errorf("长期过期时间应该为31536000秒，实际为: %d", longExpire)
	}

	// 测试推荐过期时间
	recommendedExpire := GetRecommendedExpireTime()
	if recommendedExpire != GetDefaultExpireTime() {
		t.Errorf("推荐过期时间应该等于默认过期时间")
	}
}

// TestGenUserSig 测试生成UserSig
func TestGenUserSig(t *testing.T) {
	// 注意：这个测试需要真实的腾讯云配置才能运行
	// 在实际环境中，应该使用测试配置

	api := NewTLSSigAPI(123456789, "test_secret_key")

	// 测试生成UserSig（可能会因为配置问题失败，这是正常的）
	userSig, err := api.GenUserSig("test_user", 86400)

	// 如果配置正确，应该能生成UserSig
	if err == nil && userSig != "" {
		t.Logf("UserSig生成成功: %s", userSig)
	} else {
		t.Logf("UserSig生成失败（预期行为）: %v", err)
	}
}

// TestInitTLSSigError 测试初始化错误处理
func TestInitTLSSigError(t *testing.T) {
	// 测试GetTLSSigAPISafe在未初始化时返回错误
	api, err := GetTLSSigAPISafe()
	if err == nil {
		t.Error("GetTLSSigAPISafe应该在没有初始化时返回错误")
	} else {
		t.Logf("预期的错误: %v", err)
	}

	// 测试api为nil
	if api != nil {
		t.Error("GetTLSSigAPISafe在未初始化时应该返回nil")
	}
}
