package config

import (
	"context"
	"os"
	"path/filepath"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfigParams struct {
	Endpoint   string `mapstructure:"endpoint"`
	Path       string `mapstructure:"path"`
	ServerHost string `mapstructure:"server_host"`
	ServerPort int    `mapstructure:"server_port"`
}

func EtcdConfig(endpoint string, path string) {
	var err error
	if _, err = os.Stat("./etcd"); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if os.IsNotExist(err) {
		err = os.Mkdir("./etcd", 0755)
		if err != nil {
			panic(err)
		}
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	var ctx, cancel = context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	resp, err := client.Get(ctx, path)
	if err != nil {
		panic(err)
	}

	var key = string(resp.Kvs[0].Key)
	var value = resp.Kvs[0].Value

	var basename = filepath.Base(key)
	var configPath = filepath.Join("./etcd", basename)
	err = os.WriteFile(configPath, value, 0666)
	if err != nil {
		panic(err)
	}

	InitConfig(configPath, Yaml, ".")
}

func InitEtcdConfigWithParams(cfg *EtcdConfigParams) {
	if cfg == nil {
		panic("etcd配置不能为空")
	}

	EtcdConfig(
		cfg.Endpoint,
		cfg.Path,
	)
}
