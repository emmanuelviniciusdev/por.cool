package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set up test environment variables
	testEnvVars := map[string]string{
		"MARIADB_HOST":         "testhost",
		"MARIADB_PORT":         "3307",
		"MARIADB_USER":         "testuser",
		"MARIADB_PASSWORD":     "testpass",
		"MARIADB_DATABASE":     "testdb",
		"MONGODB_URI":          "mongodb://testhost:27017",
		"MONGODB_DATABASE":     "testmongodb",
		"RABBITMQ_URI":         "amqp://guest:guest@testhost:5672/",
		"RABBITMQ_QUEUE_NAME":  "test-queue",
		"INGESTION_BATCH_SIZE": "50",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Clean up after test
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Verify MariaDB config
	if cfg.MariaDB.Host != "testhost" {
		t.Errorf("MariaDB.Host = %s, want testhost", cfg.MariaDB.Host)
	}
	if cfg.MariaDB.Port != 3307 {
		t.Errorf("MariaDB.Port = %d, want 3307", cfg.MariaDB.Port)
	}
	if cfg.MariaDB.User != "testuser" {
		t.Errorf("MariaDB.User = %s, want testuser", cfg.MariaDB.User)
	}
	if cfg.MariaDB.Password != "testpass" {
		t.Errorf("MariaDB.Password = %s, want testpass", cfg.MariaDB.Password)
	}
	if cfg.MariaDB.Database != "testdb" {
		t.Errorf("MariaDB.Database = %s, want testdb", cfg.MariaDB.Database)
	}

	// Verify MongoDB config
	if cfg.MongoDB.URI != "mongodb://testhost:27017" {
		t.Errorf("MongoDB.URI = %s, want mongodb://testhost:27017", cfg.MongoDB.URI)
	}
	if cfg.MongoDB.Database != "testmongodb" {
		t.Errorf("MongoDB.Database = %s, want testmongodb", cfg.MongoDB.Database)
	}

	// Verify RabbitMQ config
	if cfg.RabbitMQ.URI != "amqp://guest:guest@testhost:5672/" {
		t.Errorf("RabbitMQ.URI = %s, want amqp://guest:guest@testhost:5672/", cfg.RabbitMQ.URI)
	}
	if cfg.RabbitMQ.QueueName != "test-queue" {
		t.Errorf("RabbitMQ.QueueName = %s, want test-queue", cfg.RabbitMQ.QueueName)
	}

	// Verify Ingestion config
	if cfg.Ingestion.BatchSize != 50 {
		t.Errorf("Ingestion.BatchSize = %d, want 50", cfg.Ingestion.BatchSize)
	}
}

func TestLoadWithDefaults(t *testing.T) {
	// Clear all environment variables
	envVars := []string{
		"MARIADB_HOST", "MARIADB_PORT", "MARIADB_USER", "MARIADB_PASSWORD",
		"MARIADB_DATABASE", "MONGODB_URI", "MONGODB_DATABASE",
		"RABBITMQ_URI", "RABBITMQ_QUEUE_NAME",
		"INGESTION_BATCH_SIZE",
		"OPENSEARCH_ENABLED", "OPENSEARCH_URL", "OPENSEARCH_USERNAME",
		"OPENSEARCH_PASSWORD", "OPENSEARCH_INDEX_PREFIX", "OPENSEARCH_RETENTION_DAYS",
	}
	for _, key := range envVars {
		os.Unsetenv(key)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Verify defaults
	if cfg.MariaDB.Host != "localhost" {
		t.Errorf("MariaDB.Host = %s, want localhost", cfg.MariaDB.Host)
	}
	if cfg.MariaDB.Port != 3306 {
		t.Errorf("MariaDB.Port = %d, want 3306", cfg.MariaDB.Port)
	}
	if cfg.RabbitMQ.URI != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("RabbitMQ.URI = %s, want amqp://guest:guest@localhost:5672/", cfg.RabbitMQ.URI)
	}
	if cfg.RabbitMQ.QueueName != "porcool-ingestion-non-relational-database-to-relational-database" {
		t.Errorf("RabbitMQ.QueueName = %s, want porcool-ingestion-non-relational-database-to-relational-database", cfg.RabbitMQ.QueueName)
	}
	if cfg.Ingestion.BatchSize != 100 {
		t.Errorf("Ingestion.BatchSize = %d, want 100", cfg.Ingestion.BatchSize)
	}

	// Verify OpenSearch defaults
	if cfg.OpenSearch.Enabled {
		t.Error("OpenSearch.Enabled should be false by default")
	}
	if cfg.OpenSearch.URL != "" {
		t.Errorf("OpenSearch.URL = %s, want empty", cfg.OpenSearch.URL)
	}
	if cfg.OpenSearch.IndexPrefix != "porcool-ingestion-non-relational-database-to-relational-database" {
		t.Errorf("OpenSearch.IndexPrefix = %s, want porcool-ingestion-non-relational-database-to-relational-database", cfg.OpenSearch.IndexPrefix)
	}
	if cfg.OpenSearch.RetentionDays != 90 {
		t.Errorf("OpenSearch.RetentionDays = %d, want 90", cfg.OpenSearch.RetentionDays)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	os.Setenv("MARIADB_PORT", "invalid")
	defer os.Unsetenv("MARIADB_PORT")

	_, err := Load()
	if err == nil {
		t.Error("Load() should return error for invalid port")
	}
}

func TestMariaDBConfigDSN(t *testing.T) {
	cfg := MariaDBConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "secret",
		Database: "porcool",
	}

	expected := "root:secret@tcp(localhost:3306)/porcool?parseTime=true&charset=utf8mb4"
	if dsn := cfg.DSN(); dsn != expected {
		t.Errorf("DSN() = %s, want %s", dsn, expected)
	}
}

func TestGetEnv(t *testing.T) {
	// Test with set variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	if value := getEnv("TEST_VAR", "default"); value != "test_value" {
		t.Errorf("getEnv() = %s, want test_value", value)
	}

	// Test with unset variable
	if value := getEnv("UNSET_VAR", "default"); value != "default" {
		t.Errorf("getEnv() = %s, want default", value)
	}
}
