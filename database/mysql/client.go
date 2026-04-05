package mysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"

	"gorm.io/plugin/dbresolver"

	log "github.com/txze/wzkj-common/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置结构体

type Config struct {
	Name      string   // 客户端名称
	MasterDSN string   // 主数据库连接字符串
	SlaveDSNs []string // 从数据库连接字符串列表
	IdleConns int      // 最大空闲连接数
	MaxConns  int      // 最大打开连接数
}

// GormClient 数据库客户端结构体，支持主从架构
// - master: 主数据库连接（包含读写分离配置）
// - masterDSN: 主数据库连接字符串
// - slaveDSNs: 从数据库连接字符串列表
type GormClient struct {
	master    *gorm.DB
	masterDSN string
	slaveDSNs []string
}

// NewClient 创建一个新的数据库客户端实例
func NewClient() *GormClient {
	return &GormClient{
		slaveDSNs: make([]string, 0),
	}
}

// Client 全局默认数据库客户端
var Client *GormClient

// clientMap 客户端映射，用于存储多个命名客户端
var clientMap = map[string]*GormClient{}

// Dial 创建一个新的数据库客户端
// - name: 客户端名称
// - masterDSN: 主数据库连接字符串
// - slaveDSNs: 从数据库连接字符串列表
// - idle: 最大空闲连接数
// - max: 最大打开连接数
func Dial(name string, masterDSN string, slaveDSNs []string, idle, max int) *GormClient {
	Client = NewClient()
	log.Info("MYSQL 将开始连接服务器...")
	gormClient, err := Client.Dial(masterDSN, slaveDSNs, idle, max)
	if err != nil {
		panic(err)
	}

	if name == "" {
		name = "default"
	}
	clientMap[name] = gormClient

	return gormClient
}

// DialWithConfig 使用配置结构体创建一个新的数据库客户端
func DialWithConfig(config Config) *GormClient {
	return Dial(config.Name, config.MasterDSN, config.SlaveDSNs, config.IdleConns, config.MaxConns)
}

// Dial 建立数据库连接
// - masterDSN: 主数据库连接字符串
// - slaveDSNs: 从数据库连接字符串列表
// - idle: 最大空闲连接数
// - max: 最大打开连接数
func (c *GormClient) Dial(masterDSN string, slaveDSNs []string, idle, max int) (*GormClient, error) {
	var err error
	// 连接主库
	c.master, err = gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   NewGormLogger(),
	})
	if err != nil {
		return c, err
	}
	c.master.Logger.LogMode(logger.Info)
	c.master = c.master.Debug()

	// 配置连接池
	sqlMasterDB, err := c.master.DB()
	if err != nil {
		return c, err
	}

	if idle <= 0 {
		idle = 10
	}
	if max <= 0 {
		max = 50
	}
	sqlMasterDB.SetMaxIdleConns(idle)
	sqlMasterDB.SetMaxOpenConns(max)
	sqlMasterDB.SetConnMaxLifetime(time.Hour)

	// 配置从库，使用 GORM 的 dbresolver 插件实现读写分离
	if len(slaveDSNs) > 0 {
		replicas := make([]gorm.Dialector, 0, len(slaveDSNs))
		for _, dsn := range slaveDSNs {
			replicas = append(replicas, mysql.Open(dsn))
		}

		// 使用 dbresolver 插件配置读写分离，采用轮询策略
		err = c.master.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{}, // 轮询策略
		}).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(time.Hour).SetMaxIdleConns(idle).SetMaxOpenConns(max))
		if err != nil {
			log.Error("配置读写分离失败: " + err.Error())
			// 继续执行，即使读写分离配置失败，也可以使用主库
		}
	}

	c.masterDSN = masterDSN
	c.slaveDSNs = slaveDSNs

	return c, nil
}

// Master 获取主数据库连接，用于写操作
func (c *GormClient) Master() *gorm.DB {
	return c.master
}

// Slave 获取数据库连接，用于读操作
// 使用 GORM 的 dbresolver 插件时，会自动路由到从库
func (c *GormClient) Slave() *gorm.DB {
	// 当使用 dbresolver 插件时，GORM 会自动处理读写分离
	// 读操作会自动路由到从库，写操作会自动路由到主库
	return c.master
}

// GetClient 根据名称获取数据库客户端
func GetClient(name string) *GormClient {
	return clientMap[name]
}

// CheckHealth 检查数据库连接的健康状态
// - 检查主库健康状态
// - 如果主库不健康，尝试故障转移
func (c *GormClient) CheckHealth() bool {
	// 检查主库健康状态
	if c.master != nil {
		sqlDB, err := c.master.DB()
		if err != nil {
			log.Error("获取主库连接失败: " + err.Error())
			// 主库不可用时尝试故障转移
			c.HandleMasterFailover()
			return false
		}
		if err := sqlDB.Ping(); err != nil {
			log.Error("主数据库宕机: " + err.Error())
			// 主库不可用时尝试故障转移
			c.HandleMasterFailover()
			return false
		}
	}

	// 注意：使用 dbresolver 时，从库健康检查由 GORM 自动处理
	// GORM 会自动跳过不可用的从库

	return true
}

// HandleMasterFailover 处理主库故障转移
// 在 MHA + Keepalived + havip 架构下，虚拟 IP 会自动漂移到新的主库
// 这里我们重新连接到新的主库（通过相同的虚拟 IP）
// 并更新从库列表，移除已经变为主库的实例
func (c *GormClient) HandleMasterFailover() {
	log.Info("开始主库故障转移...")

	// 尝试重新连接到主库（虚拟 IP 已漂移到新主库）
	newMaster, err := gorm.Open(mysql.Open(c.masterDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   NewGormLogger(),
	})
	if err != nil {
		log.Error("重新连接主库失败: " + err.Error())
		return
	}

	// 配置新主库
	newMaster.Logger.LogMode(logger.Info)
	newMaster = newMaster.Debug()

	// 配置连接池
	sqlMasterDB, err := newMaster.DB()
	if err != nil {
		log.Error("获取新主库连接失败: " + err.Error())
		return
	}

	// 保持与原主库相同的连接池配置
	sqlMasterDB.SetMaxIdleConns(10)
	sqlMasterDB.SetMaxOpenConns(50)
	sqlMasterDB.SetConnMaxLifetime(time.Hour)

	// 更新从库列表，移除已经变为主库的实例
	// 在 MHA 架构下，当从库被提升为主库后，它应该从从库列表中移除
	updatedSlaveDSNs := make([]string, 0, len(c.slaveDSNs))
	for _, dsn := range c.slaveDSNs {
		// 检查该从库是否已经变为主库
		// 注意：在实际部署中，你可能需要根据具体的架构和配置来判断
		// 这里我们假设主库使用虚拟 IP，从库使用实际 IP，所以不需要移除
		// 如果从库提升为主库后其连接字符串会与主库相同，则需要跳过
		updatedSlaveDSNs = append(updatedSlaveDSNs, dsn)
	}

	// 重新配置从库
	if len(updatedSlaveDSNs) > 0 {
		replicas := make([]gorm.Dialector, 0, len(updatedSlaveDSNs))
		for _, dsn := range updatedSlaveDSNs {
			replicas = append(replicas, mysql.Open(dsn))
		}

		// 使用 dbresolver 插件配置读写分离
		err = newMaster.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{}, // 随机策略
		}).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(time.Hour).SetMaxIdleConns(10).SetMaxOpenConns(50))
		if err != nil {
			log.Error("配置读写分离失败: " + err.Error())
			// 继续执行，即使读写分离配置失败，也可以使用主库
		}
	}

	// 替换主库连接和从库列表
	c.master = newMaster
	c.slaveDSNs = updatedSlaveDSNs

	log.Info("主库故障转移完成: 重新连接到新主库")
}

// StartHealthChecks 启动周期性健康检查
// - interval: 健康检查间隔时间
func (c *GormClient) StartHealthChecks(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.CheckHealth()
		}
	}()
}
