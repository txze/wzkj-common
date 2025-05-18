package config

import (
	"github.com/spf13/viper"
)

const (
	Yaml = "yaml"
)

func InitConfig(name string, type_ string, paths ...string) {
	viper.SetConfigName(name)
	viper.SetConfigType(type_)

	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
