package ingestion

import (
	"testing"

	"github.com/porcool/ingestion/internal/config"
)

func TestNewService(t *testing.T) {
	cfg := &config.Config{
		Ingestion: config.IngestionConfig{
			BatchSize: 100,
		},
	}

	svc := NewService(nil, nil, cfg)

	if svc == nil {
		t.Error("NewService() returned nil")
	}

	if svc.cfg != cfg {
		t.Error("NewService() didn't set config correctly")
	}
}

func TestConfig_IngestionConfig(t *testing.T) {
	cfg := config.IngestionConfig{
		BatchSize: 50,
	}

	if cfg.BatchSize != 50 {
		t.Errorf("BatchSize = %d, want 50", cfg.BatchSize)
	}
}

func TestConfig_MariaDBConfig(t *testing.T) {
	cfg := config.MariaDBConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		Database: "porcool",
	}

	dsn := cfg.DSN()
	expected := "root:password@tcp(localhost:3306)/porcool?parseTime=true&charset=utf8mb4"

	if dsn != expected {
		t.Errorf("DSN() = %s, want %s", dsn, expected)
	}
}

func TestConfig_MongoDBConfig(t *testing.T) {
	cfg := config.MongoDBConfig{
		URI:      "mongodb://localhost:27017",
		Database: "porcool",
	}

	if cfg.URI != "mongodb://localhost:27017" {
		t.Errorf("URI = %s, want mongodb://localhost:27017", cfg.URI)
	}

	if cfg.Database != "porcool" {
		t.Errorf("Database = %s, want porcool", cfg.Database)
	}
}

func TestConfig_RabbitMQConfig(t *testing.T) {
	cfg := config.RabbitMQConfig{
		URI:       "amqp://guest:guest@localhost:5672/",
		QueueName: "porcool-ingestion-nosql-to-sql-database",
	}

	if cfg.URI != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("URI = %s, want amqp://guest:guest@localhost:5672/", cfg.URI)
	}

	if cfg.QueueName != "porcool-ingestion-nosql-to-sql-database" {
		t.Errorf("QueueName = %s, want porcool-ingestion-nosql-to-sql-database", cfg.QueueName)
	}
}

func TestExtractDocIDs(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "slice of interface with strings",
			input:    []interface{}{"doc1", "doc2", "doc3"},
			expected: []string{"doc1", "doc2", "doc3"},
		},
		{
			name:     "slice of strings",
			input:    []string{"doc1", "doc2"},
			expected: []string{"doc1", "doc2"},
		},
		{
			name:     "single string",
			input:    "doc1",
			expected: []string{"doc1"},
		},
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			input:    []interface{}{},
			expected: nil,
		},
		{
			name:     "mixed types in slice",
			input:    []interface{}{"doc1", 123, "doc2"},
			expected: []string{"doc1", "doc2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDocIDs(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("extractDocIDs() returned %d items, want %d", len(result), len(tt.expected))
				return
			}

			for i, id := range result {
				if id != tt.expected[i] {
					t.Errorf("extractDocIDs()[%d] = %s, want %s", i, id, tt.expected[i])
				}
			}
		})
	}
}

// TestFormatSpendingDate tests the formatSpendingDate function
// which converts various date formats to the standard YYYY/MM format for MariaDB
func TestFormatSpendingDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "YYYYMM format",
			input:    "202312",
			expected: "2023/12",
		},
		{
			name:     "YYYY-MM format",
			input:    "2023-12",
			expected: "2023/12",
		},
		{
			name:     "YYYY/MM format already correct",
			input:    "2023/12",
			expected: "2023/12",
		},
		{
			name:     "YYYYMM format with leading zeros",
			input:    "202301",
			expected: "2023/01",
		},
		{
			name:     "YYYY-MM format January",
			input:    "2024-01",
			expected: "2024/01",
		},
		{
			name:     "short invalid format returns as-is",
			input:    "2023",
			expected: "2023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSpendingDate(tt.input)
			if result != tt.expected {
				t.Errorf("formatSpendingDate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestParseSpendingDateToTime tests the parseSpendingDateToTime function
func TestParseSpendingDateToTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectYear  int
		expectMonth int
	}{
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "YYYY/MM format",
			input:       "2023/12",
			expectError: false,
			expectYear:  2023,
			expectMonth: 12,
		},
		{
			name:        "YYYY-MM format",
			input:       "2023-12",
			expectError: false,
			expectYear:  2023,
			expectMonth: 12,
		},
		{
			name:        "invalid format",
			input:       "2023",
			expectError: true,
		},
		{
			name:        "YYYY/MM January",
			input:       "2024/01",
			expectError: false,
			expectYear:  2024,
			expectMonth: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseSpendingDateToTime(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("parseSpendingDateToTime(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parseSpendingDateToTime(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result.Year() != tt.expectYear {
				t.Errorf("parseSpendingDateToTime(%q) year = %d, want %d", tt.input, result.Year(), tt.expectYear)
			}

			if int(result.Month()) != tt.expectMonth {
				t.Errorf("parseSpendingDateToTime(%q) month = %d, want %d", tt.input, result.Month(), tt.expectMonth)
			}

			// Day should always be 1
			if result.Day() != 1 {
				t.Errorf("parseSpendingDateToTime(%q) day = %d, want 1", tt.input, result.Day())
			}
		})
	}
}

// TestAddMonths tests the addMonths function
func TestAddMonths(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		months   int
		expected string
	}{
		{
			name:     "add 1 month",
			date:     "2023/11",
			months:   1,
			expected: "2023/12",
		},
		{
			name:     "add 1 month cross year",
			date:     "2023/12",
			months:   1,
			expected: "2024/01",
		},
		{
			name:     "add 12 months",
			date:     "2023/06",
			months:   12,
			expected: "2024/06",
		},
		{
			name:     "add 0 months",
			date:     "2023/06",
			months:   0,
			expected: "2023/06",
		},
		{
			name:     "invalid date returns empty",
			date:     "invalid",
			months:   1,
			expected: "",
		},
		{
			name:     "empty date returns empty",
			date:     "",
			months:   1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addMonths(tt.date, tt.months)
			if result != tt.expected {
				t.Errorf("addMonths(%q, %d) = %q, want %q", tt.date, tt.months, result, tt.expected)
			}
		})
	}
}

// TestGenerateMonthRange tests the generateMonthRange function
func TestGenerateMonthRange(t *testing.T) {
	tests := []struct {
		name     string
		start    string
		end      string
		expected []string
	}{
		{
			name:     "same month",
			start:    "2023/06",
			end:      "2023/06",
			expected: []string{"2023/06"},
		},
		{
			name:     "three months",
			start:    "2023/10",
			end:      "2023/12",
			expected: []string{"2023/10", "2023/11", "2023/12"},
		},
		{
			name:     "cross year",
			start:    "2023/11",
			end:      "2024/02",
			expected: []string{"2023/11", "2023/12", "2024/01", "2024/02"},
		},
		{
			name:     "full year",
			start:    "2023/01",
			end:      "2023/12",
			expected: []string{"2023/01", "2023/02", "2023/03", "2023/04", "2023/05", "2023/06", "2023/07", "2023/08", "2023/09", "2023/10", "2023/11", "2023/12"},
		},
		{
			name:     "invalid start returns nil",
			start:    "invalid",
			end:      "2023/12",
			expected: nil,
		},
		{
			name:     "invalid end returns nil",
			start:    "2023/01",
			end:      "invalid",
			expected: nil,
		},
		{
			name:     "end before start returns empty",
			start:    "2023/12",
			end:      "2023/01",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateMonthRange(tt.start, tt.end)

			if len(result) != len(tt.expected) {
				t.Errorf("generateMonthRange(%q, %q) returned %d items, want %d", tt.start, tt.end, len(result), len(tt.expected))
				t.Errorf("got: %v, want: %v", result, tt.expected)
				return
			}

			for i, month := range result {
				if month != tt.expected[i] {
					t.Errorf("generateMonthRange(%q, %q)[%d] = %q, want %q", tt.start, tt.end, i, month, tt.expected[i])
				}
			}
		})
	}
}

// TestServiceName tests the service name constant is correct
func TestServiceName(t *testing.T) {
	expectedName := "porcool-ingestion-non-relational-database-to-relational-database"
	if serviceName != expectedName {
		t.Errorf("serviceName = %q, want %q", serviceName, expectedName)
	}
}
