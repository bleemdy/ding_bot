package util

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	ParseConfig()
}

func ParseConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("app_secret", "")
	viper.SetDefault("app_key", "")
	err := viper.ReadInConfig()
	if err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatalf("write config failed: %v", err)
		}
	}
	viper.WatchConfig()
}
