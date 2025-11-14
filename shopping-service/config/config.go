package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	PaymentService PaymentServiceConfig
}
type PaymentServiceConfig struct {
	BaseURI string
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	MongoURI string
	DBName   string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("DB_NAME", "shopping_db")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		log.Println("No .env file found, using defaults")
	}

	viper.AutomaticEnv()

	var config Config
	config.Server.Port = viper.GetString("PORT")
	config.Database.MongoURI = viper.GetString("MONGO_URI")
	config.Database.DBName = viper.GetString("DB_NAME")

	config.PaymentService.BaseURI = viper.GetString("PAYMENT_SERVICE_BASE_URI")

	return &config, nil
}
