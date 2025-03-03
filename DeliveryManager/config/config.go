package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUri                  string `mapstructure:"USER_DB_URI"`
	ServiceName            string `mapstructure:"SERVICE_NAME"`
	ServicePort            string `mapstructure:"SERVICE_PORT"`
	RabbitMQURL            string `mapstructure:"RABBITMQ_URL"`
	NotificationServiceURL string `mapstructure:"NOTIFICATION_SERVICE_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	return config, err
}
