package config

import (
	"log"

	"github.com/spf13/viper"
)

type PaymentConfig struct {
	Server struct {
		Port string
	}
	Database struct {
		MongoURI string
		DBName   string
	}
}

func LoadPaymentConfig() (*PaymentConfig, error) {
	viper.SetDefault("PORT_PAYMENT", "9061")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("PAYMENT_DB_NAME", "payment_db")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	}
	cfg := &PaymentConfig{}
	cfg.Server.Port = viper.GetString("PORT_PAYMENT")
	cfg.Database.MongoURI = viper.GetString("MONGO_URI")
	cfg.Database.DBName = viper.GetString("PAYMENT_DB_NAME")
	return cfg, nil
}
