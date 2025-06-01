package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	mongoURL := os.Getenv("ME_CONFIG_MONGODB_URL")
	if mongoURL == "" {
		return fmt.Errorf("ME_CONFIG_MONGODB_URL is not set")
	}

	mongoName := os.Getenv("MONGO_DATABASE")
	if mongoName == "" {
		return fmt.Errorf("MONGO_DATABASE is not set")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return fmt.Errorf("KAFKA_BROKERS is not set")
	}

	otelEndpoint := os.Getenv("OTEL_ADDR")
	if otelEndpoint == "" {
		return fmt.Errorf("OTEL_ADDR is not set")
	}

	return nil
}

func GetMongoURL() string {
	return os.Getenv("ME_CONFIG_MONGODB_URL")
}

func GetMongoName() string {
	return os.Getenv("MONGO_DATABASE")
}

func GetPort() string {
	return GetEnvDefault("PORT", "8090")
}

func GetAppName() string {
	return GetEnvDefault("APP_NAME", "board")
}

func GetKafkaBrokers() []string {
	return strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
}

func GetOtelEndpoint() string {
	return os.Getenv("OTEL_ADDR")
}

func GetEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
