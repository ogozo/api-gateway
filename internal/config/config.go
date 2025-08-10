package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort          string `mapstructure:"HTTP_PORT"`
	UserServiceURL    string `mapstructure:"USER_SERVICE_URL"`
	ProductServiceURL string `mapstructure:"PRODUCT_SERVICE_URL"`
	CartServiceURL    string `mapstructure:"CART_SERVICE_URL"`
	OrderServiceURL   string `mapstructure:"ORDER_SERVICE_URL"`
	JWTSecretKey      string `mapstructure:"JWT_SECRET_KEY"`
}

var AppConfig *Config
func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env file not found, reading from environment variables")
	}
	
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}
}
