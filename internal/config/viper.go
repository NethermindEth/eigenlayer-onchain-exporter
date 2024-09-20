package config

import (
	"github.com/spf13/viper"
)

// GetConfig reads the configuration from the given path and returns a Config struct.
func GetConfig(configPath string) (*Config, error) {
	// Set the config file name and type
	viper.SetConfigName("eoe-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// If a config path is provided, use it
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	// Set the environment prefix and automatically use environment variables.
	// TODO: Overwriting the config file with environment variables is not working
	// as expected. We need to redefine the Config struct to make easier bindings
	// with environment variables.
	viper.SetEnvPrefix("EOE")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	return &c, viper.Unmarshal(&c)
}
