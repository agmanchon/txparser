package config

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

const (
	ConfigPath = "configPath"
)

type InstallConfig struct {
	EthereumUrl string `mapstructure:"ethereumUrl" valid:"required"`
	LogLevel    string `mapstructure:"logLevel"`
}

func ReadInstallConfig() (*InstallConfig, error) {
	// find config file
	viper.SetConfigType("yaml")
	var configPath = GetConfigPath()
	if len(configPath) < 1 {
		return nil, errors.New("configuration path is required")
	}

	// read config file
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("configuration could not be read: %w", err)
	}

	// unmarshall config file content
	var installConfig InstallConfig
	err = viper.Unmarshal(&installConfig)
	if err != nil {
		return nil, fmt.Errorf("configuration could not be unmarshalled: %w", err)
	}

	// validate config
	valid, err := validateInstallConfig(installConfig)
	if !valid {
		return nil, fmt.Errorf("configuration is not valid: %w", err)
	}

	return &installConfig, nil
}

func GetConfigPath() string {
	return viper.GetString(ConfigPath)
}

func validateInstallConfig(theConfig InstallConfig) (bool, error) {
	valid, err := govalidator.ValidateStruct(theConfig)
	if !valid || err != nil {
		return valid, err
	}
	return true, nil
}
