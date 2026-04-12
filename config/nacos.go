package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v2"
)

// ConfigFormat 配置文件格式
type ConfigFormat string

const (
	FormatJSON ConfigFormat = "json"
	FormatYAML ConfigFormat = "yaml"
	FormatYML  ConfigFormat = "yml"
)

// NacosConfig Nacos 配置中心客户端配置
// 包含客户端配置和服务器配置列表
type NacosConfig struct {
	ClientConfig  constant.ClientConfig
	ServerConfigs []constant.ServerConfig
}

var (
	client     config_client.IConfigClient
	clientOnce sync.Once
	clientLock sync.RWMutex
)

// NewNacosConfig 创建 Nacos 配置实例
// 默认配置：
// - TimeoutMs: 10000ms
// - NotLoadCacheAtStart: true
// - LogLevel: "info"
func NewNacosConfig() *NacosConfig {
	return &NacosConfig{
		ClientConfig: constant.ClientConfig{
			TimeoutMs:           10000,
			NotLoadCacheAtStart: true,
			LogLevel:            "info",
		},
	}
}

// SetAuth 设置用户名密码鉴权
func (c *NacosConfig) SetAuth(username, password string) {
	c.ClientConfig.Username = username
	c.ClientConfig.Password = password
}

// SetNamespaceId 设置命名空间 ID
func (c *NacosConfig) SetNamespaceId(namespaceId string) {
	c.ClientConfig.NamespaceId = namespaceId
}

// SetTimeoutMs 设置请求超时时间（毫秒）
func (c *NacosConfig) SetTimeoutMs(timeoutMs uint64) {
	c.ClientConfig.TimeoutMs = timeoutMs
}

// SetLogLevel 设置日志级别
// 可选值: debug, info, warn, error
func (c *NacosConfig) SetLogLevel(logLevel string) {
	c.ClientConfig.LogLevel = logLevel
}

// SetLogDir 设置日志存储目录
func (c *NacosConfig) SetLogDir(logDir string) {
	c.ClientConfig.LogDir = logDir
}

// SetCacheDir 设置缓存目录
func (c *NacosConfig) SetCacheDir(cacheDir string) {
	c.ClientConfig.CacheDir = cacheDir
}

// SetUpdateThreadNum 设置监听配置变化的并发线程数
func (c *NacosConfig) SetUpdateThreadNum(num int) {
	c.ClientConfig.UpdateThreadNum = num
}

// SetNotLoadCacheAtStart 设置启动时是否不读取缓存
func (c *NacosConfig) SetNotLoadCacheAtStart(notLoad bool) {
	c.ClientConfig.NotLoadCacheAtStart = notLoad
}

// SetUpdateCacheWhenEmpty 设置当服务返回空实例列表时是否更新缓存
func (c *NacosConfig) SetUpdateCacheWhenEmpty(update bool) {
	c.ClientConfig.UpdateCacheWhenEmpty = update
}

// AddServerConfig 添加服务器配置
func (c *NacosConfig) AddServerConfig(ipAddr string, port uint64) {
	c.ServerConfigs = append(c.ServerConfigs, *constant.NewServerConfig(ipAddr, port))
}

// AddServerConfigWithOptions 添加服务器配置（带选项）
// 支持的选项: constant.WithScheme, constant.WithContextPath
func (c *NacosConfig) AddServerConfigWithOptions(ipAddr string, port uint64, options ...constant.ServerOption) {
	c.ServerConfigs = append(c.ServerConfigs, *constant.NewServerConfig(ipAddr, port, options...))
}

// ConfigItem 配置项
// 用于批量获取多个配置文件
type ConfigItem struct {
	DataId string       // 配置 ID
	Group  string       // 配置分组
	Out    interface{}  // 绑定的目标结构体指针
	Format ConfigFormat // 配置格式: json, yaml, yml
}

// GetConfigAndBind 从 Nacos 获取配置并绑定到结构体（单例模式）
// config: Nacos 配置
// dataId: 配置 ID
// group: 配置分组
// out: 目标结构体指针，配置内容将被反序列化到该结构体
// format: 配置格式，支持 json、yaml、yml
func GetConfigAndBind(config *NacosConfig, dataId string, group string, out interface{}, format ConfigFormat) (err error) {
	client, err = getSingletonClient(config)
	if err != nil {
		return fmt.Errorf("get nacos config client failed: %w", err)
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return fmt.Errorf("get config failed: %w", err)
	}

	configFormat := FormatJSON
	if format != "" {
		configFormat = format
	}

	if err := UnmarshalContent(content, out, configFormat); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	return nil
}

// GetConfigsAndBind 批量从 Nacos 获取多个配置并绑定到结构体（单例模式）
// config: Nacos 配置
// items: 配置项列表，每个配置项包含 DataId、Group、目标结构体指针和格式
func GetConfigsAndBind(config *NacosConfig, items []ConfigItem) error {
	client, err := getSingletonClient(config)
	if err != nil {
		return fmt.Errorf("get nacos config client failed: %w", err)
	}

	for _, item := range items {
		content, err := client.GetConfig(vo.ConfigParam{
			DataId: item.DataId,
			Group:  item.Group,
		})
		if err != nil {
			return fmt.Errorf("get config %s/%s failed: %w", item.Group, item.DataId, err)
		}

		configFormat := FormatJSON
		if item.Format != "" {
			configFormat = item.Format
		}

		if err := UnmarshalContent(content, item.Out, configFormat); err != nil {
			return fmt.Errorf("unmarshal config %s/%s failed: %w", item.Group, item.DataId, err)
		}
	}

	return nil
}

// ListenConfig 监听配置变更（单例模式）
// config: Nacos 配置
// dataId: 配置 ID
// group: 配置分组
// callback: 配置变更回调函数，参数为 namespace, group, dataId, content
func ListenConfig(config *NacosConfig, dataId string, group string, callback func(namespace, group, dataId, content string)) (err error) {
	client, err = getSingletonClient(config)
	if err != nil {
		return fmt.Errorf("get nacos config client failed: %w", err)
	}

	err = client.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: callback,
	})
	if err != nil {
		return fmt.Errorf("listen config failed: %w", err)
	}

	return nil
}

// ListenConfigs 批量监听多个配置变更（单例模式）
// config: Nacos 配置
// items: 配置项列表
// callback: 配置变更回调函数
func ListenConfigs(config *NacosConfig, items []ConfigItem, callback func(namespace, group, dataId, data string)) (err error) {
	client, err = getSingletonClient(config)
	if err != nil {
		return fmt.Errorf("get nacos config client failed: %w", err)
	}

	for _, item := range items {
		err = client.ListenConfig(vo.ConfigParam{
			DataId:   item.DataId,
			Group:    item.Group,
			OnChange: callback,
		})
		if err != nil {
			return fmt.Errorf("listen config %s/%s failed: %w", item.Group, item.DataId, err)
		}
	}

	return nil
}

// UnmarshalContent 根据格式反序列化配置内容
func UnmarshalContent(content string, out interface{}, format ConfigFormat) error {
	switch format {
	case FormatYAML, FormatYML:
		return yaml.Unmarshal([]byte(content), out)
	case FormatJSON:
		return json.Unmarshal([]byte(content), out)
	default:
		return fmt.Errorf("unsupported config format: %s", format)
	}
}

// ResetClient 重置单例客户端（用于测试或重新配置）
func ResetClient() {
	clientLock.Lock()
	defer clientLock.Unlock()
	client = nil
	clientOnce = sync.Once{}
}

// getSingletonClient 获取单例 Nacos 配置客户端
// 使用 sync.Once 确保客户端只被创建一次
func getSingletonClient(config *NacosConfig) (config_client.IConfigClient, error) {
	clientLock.RLock()
	if client != nil {
		clientLock.RUnlock()
		return client, nil
	}
	clientLock.RUnlock()

	clientLock.Lock()
	defer clientLock.Unlock()

	if client != nil {
		return client, nil
	}

	var err error
	clientOnce.Do(func() {
		client, err = clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &config.ClientConfig,
				ServerConfigs: config.ServerConfigs,
			},
		)
	})

	return client, err
}
