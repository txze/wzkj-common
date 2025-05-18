package config

import (
	"context"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"path/filepath"
	"time"
)

func EtcdConfig(endpoint string, path string) {
	var err error
	// 默认指定目标文件路径
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

	// 创建viper对象
	InitConfig(configPath, Yaml, ".")
}

func InitEtcdConfig(name string, paths ...string) {
	var err error
	var conf = viper.New()
	conf.SetConfigName(name)
	conf.SetConfigType(Yaml)

	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		conf.AddConfigPath(path)
	}

	err = conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	EtcdConfig(
		conf.GetString("etcd.endpoint"),
		conf.GetString("etcd.path"),
	)

	viper.Set("server.host", conf.GetString("server.host"))
	viper.Set("server.port", conf.GetInt("server.port"))
}
