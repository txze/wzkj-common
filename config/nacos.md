# Nacos 配置中心客户端

基于 `github.com/nacos-group/nacos-sdk-go/v2` 封装的配置中心客户端，提供便捷的配置获取和监听功能。

## 功能特性

- 支持从 Nacos 获取配置并绑定到任意结构体
- 支持批量获取多个配置文件
- 支持配置变更监听
- 支持私有部署场景的用户名密码鉴权
- 提供链式调用的设置方法

## 安装

```bash
go get github.com/nacos-group/nacos-sdk-go/v2
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/txze/wzkj-common/config"
)

type AppConfig struct {
    Database struct {
        Host     string `json:"host"`
        Port     int    `json:"port"`
        Username string `json:"username"`
        Password string `json:"password"`
    } `json:"database"`
    Redis struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"redis"`
}

func main() {
    // 创建配置
    nacosCfg := config.NewNacosConfig()
    
    // 设置鉴权
    nacosCfg.SetAuth("nacos", "nacos")
    
    // 设置命名空间
    nacosCfg.SetNamespaceId("public")
    
    // 添加服务器
    nacosCfg.AddServerConfig("127.0.0.1", 8848)
    nacosCfg.AddServerConfig("127.0.0.1", 8849)
    
    // 获取单个配置
    var appConfig AppConfig
    err := config.GetConfigAndBind(nacosCfg, "app.json", "DEFAULT_GROUP", &appConfig)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Database Host: %s\n", appConfig.Database.Host)
}
```

### 批量获取配置

```go
var dbConfig DatabaseConfig
var redisConfig RedisConfig
var appConfig AppConfig

items := []config.ConfigItem{
    {DataId: "database.json", Group: "DEFAULT_GROUP", Out: &dbConfig},
    {DataId: "redis.json", Group: "DEFAULT_GROUP", Out: &redisConfig},
    {DataId: "app.json", Group: "DEFAULT_GROUP", Out: &appConfig},
}

err := config.GetConfigsAndBind(nacosCfg, items)
if err != nil {
    panic(err)
}
```

### 监听配置变更

```go
err := config.ListenConfig(nacosCfg, "app.json", "DEFAULT_GROUP", func(namespace, group, dataId, data string) {
    fmt.Printf("Config changed: %s/%s\n", group, dataId)
    fmt.Printf("New content: %s\n", data)
})
```

## API 参考

### 结构体

#### NacosConfig

```go
type NacosConfig struct {
    ClientConfig  constant.ClientConfig
    ServerConfigs []constant.ServerConfig
}
```

#### ConfigItem

```go
type ConfigItem struct {
    DataId string      // 配置 ID
    Group  string      // 配置分组
    Out    interface{} // 绑定的目标结构体指针
}
```

### 方法

#### NewNacosConfig() *NacosConfig

创建 Nacos 配置实例，包含默认配置：
- TimeoutMs: 10000ms
- NotLoadCacheAtStart: true
- LogLevel: "info"

#### SetAuth(username, password string)

设置用户名密码鉴权

#### SetNamespaceId(namespaceId string)

设置命名空间 ID

#### SetTimeoutMs(timeoutMs uint64)

设置请求超时时间（毫秒）

#### SetLogLevel(logLevel string)

设置日志级别，可选值：debug, info, warn, error

#### SetLogDir(logDir string)

设置日志存储目录

#### SetCacheDir(cacheDir string)

设置缓存目录

#### SetUpdateThreadNum(num int)

设置监听配置变化的并发线程数

#### SetNotLoadCacheAtStart(notLoad bool)

设置启动时是否不读取缓存

#### SetUpdateCacheWhenEmpty(update bool)

设置当服务返回空实例列表时是否更新缓存

#### AddServerConfig(ipAddr string, port uint64)

添加服务器配置

#### AddServerConfigWithOptions(ipAddr string, port uint64, options ...constant.ServerOption)

添加服务器配置（带选项）

#### GetConfigAndBind(config *NacosConfig, dataId string, group string, out interface{}) error

获取单个配置并绑定到结构体

#### GetConfigsAndBind(config *NacosConfig, items []ConfigItem) error

批量获取多个配置并绑定到结构体

#### ListenConfig(config *NacosConfig, dataId string, group string, callback func(namespace, group, dataId, data string)) error

监听单个配置变更

#### ListenConfigs(config *NacosConfig, items []ConfigItem, callback func(namespace, group, dataId, data string)) error

批量监听多个配置变更

## 服务器配置选项

使用 `AddServerConfigWithOptions` 时可传入的选项：

| 选项 | 说明 |
|------|------|
| `constant.WithScheme(scheme)` | 设置协议（http/https） |
| `constant.WithContextPath(path)` | 设置上下文路径 |

示例：
```go
nacosCfg.AddServerConfigWithOptions(
    "127.0.0.1",
    8848,
    constant.WithScheme("https"),
    constant.WithContextPath("/nacos"),
)
```

## 完整配置示例

```go
nacosCfg := config.NewNacosConfig()

// 基础配置
nacosCfg.SetAuth("nacos", "nacos")
nacosCfg.SetNamespaceId("my-namespace")
nacosCfg.SetTimeoutMs(5000)
nacosCfg.SetLogLevel("debug")
nacosCfg.SetLogDir("/tmp/nacos/log")
nacosCfg.SetCacheDir("/tmp/nacos/cache")
nacosCfg.SetUpdateThreadNum(30)
nacosCfg.SetNotLoadCacheAtStart(true)
nacosCfg.SetUpdateCacheWhenEmpty(false)

// 服务器配置
nacosCfg.AddServerConfig("127.0.0.1", 8848)
nacosCfg.AddServerConfig("127.0.0.1", 8849)
nacosCfg.AddServerConfigWithOptions(
    "127.0.0.1",
    8850,
    constant.WithScheme("https"),
)
```

## 注意事项

1. `GetConfigAndBind` 和 `GetConfigsAndBind` 方法要求配置内容为 JSON 格式
2. 监听配置变更时，回调函数在独立的 goroutine 中执行
3. 建议在应用启动时获取配置，运行期间通过监听机制获取更新
4. 服务器配置支持多个，客户端会自动轮询

## 错误处理

所有方法返回 `error` 类型，使用时应检查错误：

```go
err := config.GetConfigAndBind(nacosCfg, "app.json", "DEFAULT_GROUP", &appConfig)
if err != nil {
    // 处理错误，如日志记录、重试等
    log.Printf("Failed to get config: %v", err)
}
```
