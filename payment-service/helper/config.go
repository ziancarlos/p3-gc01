package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		MongoURI string
		DBName   string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables only")
	}
	cfg := &Config{}
	cfg.Server.Port = viper.GetString("PORT")
	cfg.Database.MongoURI = viper.GetString("MONGO_URI")
	cfg.Database.DBName = viper.GetString("DB_NAME")
	return cfg, nil
}
