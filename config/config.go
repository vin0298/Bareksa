package config

import (
	"github.com/spf13/viper"
	"fmt"
)

func SetupConfig(initialPath string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(initialPath)
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s", err)
	}
}
