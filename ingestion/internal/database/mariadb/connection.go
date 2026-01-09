package mariadb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/porcool/ingestion/internal/config"
	"github.com/porcool/ingestion/internal/models"
)

// Connection represents a MariaDB connection
type Connection struct {
	db  *sql.DB
	cfg config.MariaDBConfig
}

// NewConnection creates a new MariaDB connection
func NewConnection(cfg config.MariaDBConfig) (*Connection, error) {
	// First connect without database to create it if needed
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&charset=utf8mb4",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Create database if not exists
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database))
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	db.Close()

	// Connect to the actual database
	db, err = sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{db: db, cfg: cfg}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}

// DB returns the underlying database connection
func (c *Connection) DB() *sql.DB {
	return c.db
}

// RunMigrations runs all database migrations
func (c *Connection) RunMigrations() error {
	migrations := []string{
		// Domain table
		`CREATE TABLE IF NOT EXISTS domain (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			source VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_domain_type_source (type, source),
			INDEX idx_domain_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// User table
		`CREATE TABLE IF NOT EXISTS user (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255),
			email VARCHAR(255) NOT NULL UNIQUE,
			fl_admin BOOLEAN NOT NULL DEFAULT FALSE,
			monthly_income DECIMAL(15,2) NOT NULL DEFAULT 0,
			fl_payment_requested BOOLEAN NOT NULL DEFAULT FALSE,
			fl_payment_pending BOOLEAN NOT NULL DEFAULT FALSE,
			fl_payment_paid BOOLEAN NOT NULL DEFAULT FALSE,
			current_spending_date VARCHAR(7),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_user_email (email),
			INDEX idx_user_source_id (source_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Financial institution table
		`CREATE TABLE IF NOT EXISTS financial_institution (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			user_id BIGINT NOT NULL,
			name VARCHAR(255) NOT NULL,
			fl_credit_card BOOLEAN NOT NULL DEFAULT FALSE,
			fl_money_movement BOOLEAN NOT NULL DEFAULT FALSE,
			fl_investment BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_fi_user_id (user_id),
			INDEX idx_fi_source_id (source_id),
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Expense automatic workflow table
		`CREATE TABLE IF NOT EXISTS expense_automatic_workflow (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			user_id BIGINT NOT NULL,
			base64_image LONGTEXT,
			description TEXT,
			extracted_expense_content_from_image LONGTEXT,
			spending_date__YYYY_MM VARCHAR(7),
			sync_processed_date TIMESTAMP NULL,
			id_sync_status BIGINT,
			processing_message TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_eaw_user_id (user_id),
			INDEX idx_eaw_sync_status (id_sync_status),
			INDEX idx_eaw_source_id (source_id),
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
			FOREIGN KEY (id_sync_status) REFERENCES domain(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Expense automatic workflow pre-saved description table
		`CREATE TABLE IF NOT EXISTS expense_automatic_workflow_pre_saved_description (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			user_id BIGINT NOT NULL,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_eawpsd_user_id (user_id),
			INDEX idx_eawpsd_source_id (source_id),
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Expense table
		`CREATE TABLE IF NOT EXISTS expense (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			user_id BIGINT NOT NULL,
			spending_date__YYYY_MM VARCHAR(7) NOT NULL,
			id_status BIGINT,
			id_type BIGINT,
			validity_period_date DATE,
			fl_indeterminate_validity_period_date BOOLEAN NOT NULL DEFAULT FALSE,
			name VARCHAR(255) NOT NULL,
			total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
			total_paid_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_expense_user_id (user_id),
			INDEX idx_expense_spending_date (spending_date__YYYY_MM),
			INDEX idx_expense_status (id_status),
			INDEX idx_expense_type (id_type),
			INDEX idx_expense_source_id (source_id),
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
			FOREIGN KEY (id_status) REFERENCES domain(id),
			FOREIGN KEY (id_type) REFERENCES domain(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Expense installment table
		`CREATE TABLE IF NOT EXISTS expense_installment (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			expense_id BIGINT NOT NULL,
			amount DECIMAL(15,2) NOT NULL DEFAULT 0,
			paid_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
			id_status BIGINT,
			due_date DATE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_ei_expense_id (expense_id),
			INDEX idx_ei_status (id_status),
			INDEX idx_ei_due_date (due_date),
			FOREIGN KEY (expense_id) REFERENCES expense(id) ON DELETE CASCADE,
			FOREIGN KEY (id_status) REFERENCES domain(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Additional balance table
		`CREATE TABLE IF NOT EXISTS additional_balance (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			user_id BIGINT NOT NULL,
			spending_date__YYYY_MM VARCHAR(7) NOT NULL,
			amount DECIMAL(15,2) NOT NULL DEFAULT 0,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_ab_user_id (user_id),
			INDEX idx_ab_spending_date (spending_date__YYYY_MM),
			INDEX idx_ab_source_id (source_id),
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Balance history table
		`CREATE TABLE IF NOT EXISTS balance_history (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			user_id BIGINT NOT NULL,
			spending_date__YYYY_MM VARCHAR(7) NOT NULL,
			amount DECIMAL(15,2) NOT NULL DEFAULT 0,
			last_month_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
			monthly_income DECIMAL(15,2) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_bh_user_id (user_id),
			INDEX idx_bh_spending_date (spending_date__YYYY_MM),
			INDEX idx_bh_source_id (source_id),
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Service payment table
		`CREATE TABLE IF NOT EXISTS service_payment (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			user_id BIGINT NOT NULL,
			service_payment_date DATE NOT NULL,
			service_payment_type_id BIGINT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_sp_user_id (user_id),
			INDEX idx_sp_payment_date (service_payment_date),
			INDEX idx_sp_source_id (source_id),
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
			FOREIGN KEY (service_payment_type_id) REFERENCES domain(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// System settings table
		`CREATE TABLE IF NOT EXISTS system_settings (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			guid VARCHAR(36) NOT NULL UNIQUE,
			source_id VARCHAR(255) NOT NULL,
			fl_block_user_registration BOOLEAN NOT NULL DEFAULT FALSE,
			fl_maintenance BOOLEAN NOT NULL DEFAULT FALSE,
			json_sync_metadata JSON,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP NULL,
			updated_by VARCHAR(255),
			INDEX idx_ss_source_id (source_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for _, migration := range migrations {
		if _, err := c.db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w\nSQL: %s", err, migration)
		}
	}

	return nil
}

// SeedDomains seeds the domain table with initial values
func (c *Connection) SeedDomains() error {
	seeds := models.GetDomainSeeds()

	for _, seed := range seeds {
		for _, name := range seed.Names {
			// Check if domain already exists
			var count int
			err := c.db.QueryRow(
				"SELECT COUNT(*) FROM domain WHERE name = ? AND type = ? AND source = ?",
				name, seed.Type, seed.Source,
			).Scan(&count)
			if err != nil {
				return fmt.Errorf("failed to check domain existence: %w", err)
			}

			// Insert if not exists
			if count == 0 {
				guid := uuid.New().String()
				_, err := c.db.Exec(
					"INSERT INTO domain (guid, name, type, source, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?)",
					guid, name, seed.Type, seed.Source, time.Now(), "system",
				)
				if err != nil {
					return fmt.Errorf("failed to insert domain: %w", err)
				}
			}
		}
	}

	return nil
}

// GetDomainID returns the ID of a domain by name, type, and source
func (c *Connection) GetDomainID(name, domainType, source string) (int64, error) {
	var id int64
	err := c.db.QueryRow(
		"SELECT id FROM domain WHERE name = ? AND type = ? AND source = ?",
		name, domainType, source,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("domain not found: %w", err)
	}
	return id, nil
}
