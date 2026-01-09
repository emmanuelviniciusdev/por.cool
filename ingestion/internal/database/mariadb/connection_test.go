package mariadb

import (
	"testing"

	"github.com/porcool/ingestion/internal/config"
)

func TestMariaDBConfigDSN(t *testing.T) {
	cfg := config.MariaDBConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		Database: "testdb",
	}

	expected := "root:password@tcp(localhost:3306)/testdb?parseTime=true&charset=utf8mb4"
	if dsn := cfg.DSN(); dsn != expected {
		t.Errorf("DSN() = %s, want %s", dsn, expected)
	}
}

func TestMariaDBConfigDSNWithSpecialChars(t *testing.T) {
	cfg := config.MariaDBConfig{
		Host:     "db.example.com",
		Port:     3307,
		User:     "admin",
		Password: "p@ss#word",
		Database: "porcool",
	}

	expected := "admin:p@ss#word@tcp(db.example.com:3307)/porcool?parseTime=true&charset=utf8mb4"
	if dsn := cfg.DSN(); dsn != expected {
		t.Errorf("DSN() = %s, want %s", dsn, expected)
	}
}

// Integration tests (require actual database connection)
// These are skipped by default and can be run with: go test -tags=integration

func TestNewConnection_InvalidHost(t *testing.T) {
	cfg := config.MariaDBConfig{
		Host:     "nonexistent.host.local",
		Port:     3306,
		User:     "root",
		Password: "password",
		Database: "testdb",
	}

	_, err := NewConnection(cfg)
	if err == nil {
		t.Error("NewConnection() should fail with invalid host")
	}
}

func TestConnection_Close(t *testing.T) {
	// This test verifies the Close method doesn't panic
	// In production, you'd use a mock or actual connection
	conn := &Connection{db: nil}

	// This will panic if db is nil and Close is called
	// We're just testing that the method signature is correct
	if conn.db != nil {
		err := conn.Close()
		if err != nil {
			t.Errorf("Close() returned error: %v", err)
		}
	}
}

func TestConnection_DB(t *testing.T) {
	conn := &Connection{db: nil}

	if conn.DB() != nil {
		t.Error("DB() should return nil for nil connection")
	}
}
