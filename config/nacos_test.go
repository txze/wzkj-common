package config

import (
	"testing"

	"github.com/hzxiao/goutil"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

func TestFullConfigExample(t *testing.T) {
	cfg := NewNacosConfig()

	cfg.SetAuth("111", "111")
	cfg.SetNamespaceId("")
	cfg.SetTimeoutMs(5000)
	cfg.SetLogLevel("debug")
	cfg.SetLogDir("./log")
	cfg.SetCacheDir("./cache")
	cfg.SetUpdateThreadNum(30)
	cfg.SetNotLoadCacheAtStart(true)
	cfg.SetUpdateCacheWhenEmpty(false)

	cfg.AddServerConfigWithOptions(
		"IP",
		8848,
		constant.WithScheme("http"),
		constant.WithContextPath("/nacos"),
	)

	if cfg.ClientConfig.Username != "admin" {
		t.Errorf("Expected Username 'admin', got '%s'", cfg.ClientConfig.Username)
	}
	if cfg.ClientConfig.Password != "password123" {
		t.Errorf("Expected Password 'password123', got '%s'", cfg.ClientConfig.Password)
	}
	if cfg.ClientConfig.NamespaceId != "my-namespace" {
		t.Errorf("Expected NamespaceId 'my-namespace', got '%s'", cfg.ClientConfig.NamespaceId)
	}
	if cfg.ClientConfig.TimeoutMs != 5000 {
		t.Errorf("Expected TimeoutMs 5000, got %d", cfg.ClientConfig.TimeoutMs)
	}
	if cfg.ClientConfig.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.ClientConfig.LogLevel)
	}

	if cfg.ClientConfig.UpdateThreadNum != 30 {
		t.Errorf("Expected UpdateThreadNum 30, got %d", cfg.ClientConfig.UpdateThreadNum)
	}
	if !cfg.ClientConfig.NotLoadCacheAtStart {
		t.Error("Expected NotLoadCacheAtStart to be true")
	}
	if cfg.ClientConfig.UpdateCacheWhenEmpty {
		t.Error("Expected UpdateCacheWhenEmpty to be false")
	}

	if cfg.ServerConfigs[0].Scheme != "http" {
		t.Errorf("Expected ServerConfigs[2].Scheme 'https', got '%s'", cfg.ServerConfigs[2].Scheme)
	}
	if cfg.ServerConfigs[0].ContextPath != "/nacos" {
		t.Errorf("Expected ServerConfigs[2].ContextPath '/nacos', got '%s'", cfg.ServerConfigs[2].ContextPath)
	}

	var data = make(goutil.Map)
	err := GetConfigAndBind(cfg, "jinzhao-dev", "DEFAULT_GROUP", &data, FormatYAML)
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err.Error())
	}

	t.Logf("Config data: %v", data)

	err = ListenConfig(cfg, "jinzhao-dev", "DEFAULT_GROUP", func(namespace, group, dataId, content string) {
		t.Logf("Received config update: %s", content)
		// 处理配置更新
		// ...
		// 解析配置内容
		if err := UnmarshalContent(content, &data, FormatYAML); err != nil {
			t.Errorf("Failed to parse config: %v", err)
			return
		}
		t.Logf("Parsed config: %v", data)
	})

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err.Error())
	}

	select {}

}
