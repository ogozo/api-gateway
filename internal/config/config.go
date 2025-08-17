package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GatewayConfig struct {
	HTTPPort             string `mapstructure:"HTTP_PORT"`
	UserServiceURL       string `mapstructure:"USER_SERVICE_URL"`
	ProductServiceURL    string `mapstructure:"PRODUCT_SERVICE_URL"`
	CartServiceURL       string `mapstructure:"CART_SERVICE_URL"`
	OrderServiceURL      string `mapstructure:"ORDER_SERVICE_URL"`
	JWTSecretKey         string `mapstructure:"JWT_SECRET_KEY"`
	OtelExporterEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtelServiceName      string `mapstructure:"OTEL_SERVICE_NAME"`
}

func LoadConfig(cfg any) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		tempLogger, _ := zap.NewProduction()
		defer tempLogger.Sync()
		tempLogger.Warn(".env file not found, reading from environment variables")
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		tempLogger, _ := zap.NewProduction()
		defer tempLogger.Sync()
		tempLogger.Fatal("Unable to decode config into struct", zap.Error(err))
	}
}
