package config

import (
	"errors"
	"log"
	"os"

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
	// Set defaults first
	viper.SetDefault("PORT_SHOPPING", "9051")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("SHOPPING_DB_NAME", "shopping_db")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Try to read .env file (optional - don't fail if not found)
	viper.SetConfigFile("../.env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		// Check if it's a file not found error
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No .env file found, using environment variables and defaults")
		} else if errors.Is(err, os.ErrNotExist) {
			log.Println("No .env file found, using environment variables and defaults")
		} else {
			// For other errors, log but don't fail (use env vars instead)
			log.Printf("Warning: Error reading .env file: %v, using environment variables and defaults", err)
		}
	}

	var config Config
	config.Server.Port = viper.GetString("PORT_SHOPPING")
	config.Database.MongoURI = viper.GetString("MONGO_URI")
	config.Database.DBName = viper.GetString("SHOPPING_DB_NAME")
	config.PaymentService.BaseURI = viper.GetString("PAYMENT_SERVICE_BASE_URI")

	return &config, nil
}
