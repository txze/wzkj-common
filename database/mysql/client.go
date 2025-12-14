package mysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"

	log "github.com/txze/wzkj-common/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormClient struct {
	master *gorm.DB
	slave  *gorm.DB
}

func NewClient() *GormClient {
	return &GormClient{}
}

var Client *GormClient
var clientMap = map[string]*GormClient{}

func Dial(name string, dialect string) *GormClient {
	Client = NewClient()
	log.Info("MYSQL will start dial server...")
	gormClient, err := Client.Dial(dialect)
	if err != nil {
		panic(err)
	}

	if name == "" {
		name = "default"
	}
	clientMap[name] = gormClient

	return gormClient
}

func (c *GormClient) Dial(dialect string) (*GormClient, error) {
	var err error
	// master dial
	c.master, err = gorm.Open(mysql.Open(dialect), &gorm.Config{
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

	sqlMasterDB.SetMaxIdleConns(5)
	sqlMasterDB.SetMaxOpenConns(10)
	sqlMasterDB.SetConnMaxLifetime(time.Hour)

	// slave dial
	c.slave, err = gorm.Open(mysql.Open(dialect), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   NewGormLogger(),
	})
	if err != nil {
		return c, err
	}
	c.slave.Logger.LogMode(logger.Info)
	c.slave = c.slave.Debug()
	sqlSlaveDB, err := c.slave.DB()
	if err != nil {
		return c, err
	}

	sqlSlaveDB.SetMaxIdleConns(10)
	sqlSlaveDB.SetMaxOpenConns(30)
	sqlSlaveDB.SetConnMaxLifetime(time.Hour)

	return c, nil
}

func (c *GormClient) Master() *gorm.DB {
	return c.master
}

func (c *GormClient) Slave() *gorm.DB {
	return c.slave
}

func GetClient(name string) *GormClient {
	return clientMap[name]
}
