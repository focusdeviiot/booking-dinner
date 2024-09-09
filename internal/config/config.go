package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Logger LoggerConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type LoggerConfig struct {
	Production bool
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config host and port
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %s", err)
	}

	return &config, nil
}

// validateConfig checks if the config is valid
func validateConfig(config *Config) error {
	if config.Server.Port == 0 {
		return fmt.Errorf("server port is required")
	}
	if config.Server.Host == "" {
		return fmt.Errorf("database host is required")
	}
	return nil
}
