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

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGODB_URI is not set")
	}

	mongoName := os.Getenv("MONGODB_NAME")
	if mongoName == "" {
		return fmt.Errorf("MONGODB_NAME is not set")
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

func GetMongoURI() string {
	return os.Getenv("MONGODB_URI")
}

func GetMongoName() string {
	return os.Getenv("MONGODB_NAME")
}

func GetPort() string {
	return GetEnvDefault("PORT", "9090")
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
