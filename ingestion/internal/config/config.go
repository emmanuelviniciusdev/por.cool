package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the ingestion service
type Config struct {
	MariaDB   MariaDBConfig
	MongoDB   MongoDBConfig
	RabbitMQ  RabbitMQConfig
	Ingestion IngestionConfig
}

// MariaDBConfig holds MariaDB connection configuration
type MariaDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// MongoDBConfig holds MongoDB connection configuration
type MongoDBConfig struct {
	URI      string
	Database string
}

// RabbitMQConfig holds RabbitMQ connection configuration
type RabbitMQConfig struct {
	URI       string
	QueueName string
}

// IngestionConfig holds ingestion process configuration
type IngestionConfig struct {
	BatchSize int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	mariaPort, err := strconv.Atoi(getEnv("MARIADB_PORT", "3306"))
	if err != nil {
		return nil, fmt.Errorf("invalid MARIADB_PORT: %w", err)
	}

	batchSize, err := strconv.Atoi(getEnv("INGESTION_BATCH_SIZE", "100"))
	if err != nil {
		return nil, fmt.Errorf("invalid INGESTION_BATCH_SIZE: %w", err)
	}

	return &Config{
		MariaDB: MariaDBConfig{
			Host:     getEnv("MARIADB_HOST", "localhost"),
			Port:     mariaPort,
			User:     getEnv("MARIADB_USER", "root"),
			Password: getEnv("MARIADB_PASSWORD", ""),
			Database: getEnv("MARIADB_DATABASE", "porcool"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "porcool"),
		},
		RabbitMQ: RabbitMQConfig{
			URI:       getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
			QueueName: getEnv("RABBITMQ_QUEUE_NAME", "porcool-ingestion-non-relational-database-to-relational-database"),
		},
		Ingestion: IngestionConfig{
			BatchSize: batchSize,
		},
	}, nil
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DSN returns the MariaDB Data Source Name
func (c *MariaDBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		c.User, c.Password, c.Host, c.Port, c.Database)
}
