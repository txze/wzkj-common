package postgres

import (
	"gorm.io/driver/postgres"
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

func Dial(dialect string) (*GormClient, error) {
	Client = NewClient()
	return Client.Dial(dialect)
}

func (c *GormClient) Dial(dialect string) (*GormClient, error) {
	var err error
	// master dial
	c.master, err = gorm.Open(postgres.Open(dialect), &gorm.Config{})
	if err != nil {
		return c, err
	}
	c.master.Logger.LogMode(logger.Info)
	c.master = c.master.Debug()

	// slave dial
	c.slave, err = gorm.Open(postgres.Open(dialect), &gorm.Config{})
	if err != nil {
		return c, err
	}
	c.slave.Logger.LogMode(logger.Info)
	c.slave = c.slave.Debug()
	return c, nil
}

func (c *GormClient) Master() *gorm.DB {
	return c.master
}

func (c *GormClient) Slave() *gorm.DB {
	return c.slave
}
