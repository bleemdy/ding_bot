package util

import (
	"github.com/spf13/viper"
	"log"
)

func ParseConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("appSecret", "")
	err := viper.ReadInConfig()
	viper.WatchConfig()
	if err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatal("write config failed: %v", err)
		}
	}
}
