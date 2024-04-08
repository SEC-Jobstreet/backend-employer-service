package utils

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The value are read by viper from a config file or environment variables.
type Config struct {
	Environment          string `mapstructure:"ENVIRONMENT"`
	RESTfulServerAddress string `mapstructure:"RESTfulServerAddress"`
	DBSource             string `mapstructure:"DB_SOURCE"`
	MigrationURL         string `mapstructure:"MIGRATION_URL"`
	GRPCServerAddress    string `mapstructure:"GRPCServerAddress"`
	HTTPServerAddress    string `mapstructure:"HTTPServerAddress"`
}

// LoadConfig reads configuration from file or environment variable.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
