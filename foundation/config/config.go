// Package config provides a configuration struct and a function to read the config file
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is a struct that holds the configuration values
type Config struct {
	LogLevel    string
	ServiceName string
	MasterNode  string
}

// NewConfig reads the config file and returns a Config struct
func NewConfig() (*Config, error) {
	// Set default values for your configuration
	viper.SetDefault("LogLevel", "debug")
	viper.SetDefault("ServiceName", "assignment")
	viper.SetDefault("MasterNode", ":8080")

	// Read the config file (if it exists)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./service/config")

	err := viper.ReadInConfig()
	if err != nil {
		// Handle the error if the config file doesn't exist (you can choose to create it or use only default values)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}
	}

	// Bind environment variables to config keys
	viper.AutomaticEnv()

	// Unmarshal the config into a struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}
