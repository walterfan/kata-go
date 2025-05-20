package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

var configInitialized = false

// InitConfig initializes viper configuration only once
func InitConfig() error {
	if configInitialized {
		return nil
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	configInitialized = true
	return nil
}
