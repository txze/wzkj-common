package mysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"

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
// - master: 主数据库连接
// - slaves: 从数据库连接列表
// - masterDSN: 主数据库连接字符串
// - slaveDSNs: 从数据库连接字符串列表
// - currentSlaveIndex: 当前使用的从库索引，用于轮询负载均衡
type GormClient struct {
	master            *gorm.DB
	slaves            []*gorm.DB
	masterDSN         string
	slaveDSNs         []string
	currentSlaveIndex int
}

// NewClient 创建一个新的数据库客户端实例
func NewClient() *GormClient {
	return &GormClient{
		slaves:    make([]*gorm.DB, 0),
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
	// master dial
	c.master, err = gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   NewGormLogger(),
	})
	if err != nil {
		return c, err
	}
	c.master.Logger.LogMode(logger.Info)
	c.master = c.master.Debug()

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

	// 连接从库
	for _, dsn := range slaveDSNs {
		slave, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   NewGormLogger(),
		})
		if err != nil {
			log.Error("连接从库失败: " + err.Error())
			continue
		}
		slave.Logger.LogMode(logger.Info)
		slave = slave.Debug()
		sqlSlaveDB, err := slave.DB()
		if err != nil {
			log.Error("获取从库连接失败: " + err.Error())
			continue
		}

		sqlSlaveDB.SetMaxIdleConns(idle)
		sqlSlaveDB.SetMaxOpenConns(max)
		sqlSlaveDB.SetConnMaxLifetime(time.Hour)

		c.slaves = append(c.slaves, slave)
	}

	c.masterDSN = masterDSN
	c.slaveDSNs = slaveDSNs

	return c, nil
}

// Master 获取主数据库连接，用于写操作
func (c *GormClient) Master() *gorm.DB {
	return c.master
}

// Slave 获取从数据库连接，用于读操作
// - 实现了轮询负载均衡
// - 如果没有从库可用，会回退到主库
func (c *GormClient) Slave() *gorm.DB {
	if len(c.slaves) == 0 {
		// 没有从库可用时回退到主库
		return c.master
	}

	// 轮询负载均衡
	slave := c.slaves[c.currentSlaveIndex]
	c.currentSlaveIndex = (c.currentSlaveIndex + 1) % len(c.slaves)
	return slave
}

// GetClient 根据名称获取数据库客户端
func GetClient(name string) *GormClient {
	return clientMap[name]
}

// CheckHealth 检查所有数据库连接的健康状态
// - 检查主库健康状态
// - 检查从库健康状态
// - 如果主库不健康，尝试故障转移
// - 如果从库不健康，从列表中移除
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

	// 检查从库健康状态
	if len(c.slaves) > 0 {
		for i, slave := range c.slaves {
			sqlDB, err := slave.DB()
			if err != nil {
				log.Error("获取从库连接失败: " + err.Error())
				// 从列表中移除宕机的从库
				c.slaves = append(c.slaves[:i], c.slaves[i+1:]...)
				continue
			}
			if err := sqlDB.Ping(); err != nil {
				log.Error("从数据库宕机: " + err.Error())
				// 从列表中移除宕机的从库
				c.slaves = append(c.slaves[:i], c.slaves[i+1:]...)
			}
		}
	}

	return true
}

// HandleMasterFailover 处理主库故障转移，将一个从库提升为主库
func (c *GormClient) HandleMasterFailover() {
	if len(c.slaves) == 0 {
		log.Error("没有可用的从库进行故障转移")
		return
	}

	// 将第一个可用的从库提升为主库
	newMaster := c.slaves[0]
	c.master = newMaster

	// 从从库列表中移除被提升的从库
	c.slaves = c.slaves[1:]

	// 更新主库连接字符串（假设从库连接字符串顺序与从库列表一致）
	if len(c.slaveDSNs) > 0 {
		c.masterDSN = c.slaveDSNs[0]
		c.slaveDSNs = c.slaveDSNs[1:]
	}

	log.Info("主库故障转移完成: 将一个从库提升为主库")
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
