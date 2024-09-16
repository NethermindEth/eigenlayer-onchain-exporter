package config

import "github.com/spf13/viper"

// GetConfig reads the configuration from the given path and returns a Config struct.
func GetConfig(configPath string) (*Config, error) {
	var c Config
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &c, viper.Unmarshal(&c)
}
